package repo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"gitlab.com/makeos/mosdef/types/core"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/packfile"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
)

// packedReferenceObject represent references added to a pack file
type packedReferenceObject struct {
	name    string
	oldHash string
	newHash string
}

// packedReferences represents a collection of packed references
type packedReferences []*packedReferenceObject

// names return the names of the references
func (p *packedReferences) names() (refs []string) {
	for _, p := range *p {
		refs = append(refs, p.name)
	}
	return
}

type packObject struct {
	Type plumbing.ObjectType
	Hash plumbing.Hash
}

// objObserver implements packfile.Observer
type objectObserver struct {
	objects []*packObject
}

func (o *objectObserver) OnInflatedObjectHeader(t plumbing.ObjectType, objSize int64,
	pos int64) error {
	o.objects = append(o.objects, &packObject{Type: t})
	return nil
}

func (o *objectObserver) OnInflatedObjectContent(h plumbing.Hash, pos int64,
	crc uint32, content []byte) error {
	o.objects[len(o.objects)-1].Hash = h
	return nil
}

func (o *objectObserver) OnHeader(count uint32) error    { return nil }
func (o *objectObserver) OnFooter(h plumbing.Hash) error { return nil }

// PushReader inspects push data from git client, extracting data such as the
// pushed references, objects and object to reference mapping. It also pipes the
// pushed stream to a destination (git-receive-pack) when finished.
type PushReader struct {
	dst         io.WriteCloser
	packFile    *os.File
	buf         []byte
	references  packedReferences
	objects     []*packObject
	objectsRefs objRefMap
	repo        core.BareRepo
	updateReqCB func(ur *packp.ReferenceUpdateRequest) error
}

// newPushReader creates an instance of PushReader, and after inspection, the
// written content will be copied to dst.
func newPushReader(dst io.WriteCloser, repo core.BareRepo) (*PushReader, error) {
	packFile, err := ioutil.TempFile(os.TempDir(), "pack")
	if err != nil {
		return nil, err
	}

	return &PushReader{
		dst:         dst,
		packFile:    packFile,
		repo:        repo,
		objectsRefs: make(map[string][]string),
		objects:     []*packObject{},
		references:  []*packedReferenceObject{},
	}, nil
}

// Write implements the io.Writer interface.
func (r *PushReader) Write(p []byte) (int, error) {
	return r.packFile.Write(p)
}

// OnReferenceUpdateRequestRead sets a callback that is called after the
// push requested has been decoded but yet to be written to git.
// If the callback returns an error, the push request is aborted.
func (r *PushReader) OnReferenceUpdateRequestRead(cb func(ur *packp.ReferenceUpdateRequest) error) {
	r.updateReqCB = cb
}

// Read reads the packfile, extracting object and reference information
// and finally writes the read data to a provided destination
func (r *PushReader) Read() error {

	var err error

	// Seek to the beginning of the packfile
	r.packFile.Seek(0, 0)

	// Decode the packfile into a ReferenceUpdateRequest
	ur := packp.NewReferenceUpdateRequest()
	if err = ur.Decode(r.packFile); err != nil {
		return err
	}

	// Call OnReferenceUpdateRequestRead callback method
	if r.updateReqCB != nil {
		if err = r.updateReqCB(ur); err != nil {
			return err
		}
	}

	// Extract references from the packfile
	r.references = append(r.references, r.getReferences(ur)...)

	// Scan the packfile and extract objects hashes.
	// Confirm if the next 4 bytes are indeed 'PACK', otherwise, the packfile is invalid
	packSig := make([]byte, 4)
	r.packFile.Read(packSig)
	if string(packSig) != "PACK" {
		return r.done()
	}
	r.packFile.Seek(-4, 1)

	// Read the packfile
	scn := packfile.NewScanner(r.packFile)
	defer scn.Close()
	r.objects, err = r.getObjects(scn)
	if err != nil {
		return errors.Wrap(err, "failed to get objects")
	}

	return r.done()
}

// getObjects returns a list of objects in the packfile
func (r *PushReader) getObjects(scanner *packfile.Scanner) (objs []*packObject, err error) {
	objObserver := &objectObserver{}
	packfileParser, err := packfile.NewParserWithStorage(scanner, r.repo.GetHost(), objObserver)
	if err != nil {
		return nil, err
	}
	if _, err := packfileParser.Parse(); err != nil {
		return nil, err
	}
	return objObserver.objects, nil
}

// getReferences returns the references found in the pack buffer
func (r *PushReader) getReferences(ur *packp.ReferenceUpdateRequest) (references []*packedReferenceObject) {
	for _, cmd := range ur.Commands {
		refObj := &packedReferenceObject{
			name:    cmd.Name.String(),
			oldHash: cmd.Old.String(),
			newHash: cmd.New.String(),
		}
		references = append(references, refObj)
	}
	return
}

// done copies the written content from the inspector to dst and closes the
// destination and source readers and creates a mapping of objects to references.
func (r *PushReader) done() (err error) {

	r.packFile.Seek(0, 0)
	if _, err = io.Copy(r.dst, r.packFile); err != nil {
		return
	}

	if err = r.packFile.Close(); err != nil {
		return
	}

	if err = r.dst.Close(); err != nil {
		return
	}

	// Give git some time to process the input
	time.Sleep(100 * time.Millisecond)

	r.objectsRefs, err = r.mapObjectsToRef()
	if err != nil {
		return errors.Wrap(err, "failed to map objects to references")
	}

	os.Remove(r.packFile.Name())

	return
}

// objRefMap maps objects to the references they belong to.
type objRefMap map[string][]string

// removeRef removes a reference from the list of references an object belongs to
func (m *objRefMap) removeRef(objHash, ref string) error {
	refs, ok := (*m)[objHash]
	if !ok {
		return fmt.Errorf("object not found")
	}
	newRefs := []string{}
	for _, r := range refs {
		if r != ref {
			newRefs = append(newRefs, r)
		}
	}
	(*m)[objHash] = newRefs
	return nil
}

// getObjects returns a list of objects that map to the given ref
func (m *objRefMap) getObjectsOf(ref string) (objs []string) {
	for obj, refs := range *m {
		if funk.ContainsString(refs, ref) {
			objs = append(objs, obj)
		}
	}
	return
}

// mapObjectsToRef returns a map that pairs pushed objects to one or more
// repository references they belong to.
func (r *PushReader) mapObjectsToRef() (objRefMap, error) {
	var mappings = make(map[string][]string)

	if len(r.objects) == 0 {
		return mappings, nil
	}

	for _, ref := range r.references.names() {
		var entries []string
		var err error

		refObj, err := r.repo.Reference(plumbing.ReferenceName(ref), true)
		if err != nil {
			return nil, err
		}

		obj, err := r.repo.Object(plumbing.AnyObject, refObj.Hash())
		if err != nil {
			return nil, err
		}

		objType := obj.Type()

		if objType == plumbing.CommitObject {
			entries, err = getCommitHistory(r.repo, obj.(*object.Commit), "")
			if err != nil {
				return nil, err
			}
		}

		if objType == plumbing.TagObject {
			commit, err := obj.(*object.Tag).Commit()
			if err != nil {
				return nil, err
			}
			entries, err = getCommitHistory(r.repo, commit, "")
			if err != nil {
				return nil, err
			}
			entries = append(entries, obj.(*object.Tag).ID().String())
		}

		for _, obj := range r.objects {
			if funk.ContainsString(entries, obj.Hash.String()) {
				objRefs, ok := mappings[obj.Hash.String()]
				if !ok {
					objRefs = []string{}
				}
				if !funk.ContainsString(objRefs, ref) {
					objRefs = append(objRefs, ref)
				}
				mappings[obj.Hash.String()] = objRefs

			}
		}
	}

	return mappings, nil
}
