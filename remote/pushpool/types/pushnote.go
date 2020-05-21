package types

import (
	"time"

	"github.com/vmihailenco/msgpack"
	"gitlab.com/makeos/mosdef/remote/types"
	"gitlab.com/makeos/mosdef/types/state"
	"gitlab.com/makeos/mosdef/util"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/storage"
)

type PushNotice interface {
	GetTargetRepo() LocalRepo
	SetTargetRepo(repo LocalRepo)
	GetPusherKeyID() []byte
	GetPusherAddress() util.Address
	GetPusherAccountNonce() uint64
	GetPusherKeyIDString() string
	EncodeMsgpack(enc *msgpack.Encoder) error
	DecodeMsgpack(dec *msgpack.Decoder) error
	Bytes(recompute ...bool) []byte
	BytesNoCache() []byte
	BytesNoSig() []byte
	GetPushedObjects() []string
	GetEcoSize() uint64
	GetNodePubKey() util.Bytes32
	GetNodeSignature() []byte
	GetRepoName() string
	GetNamespace() string
	GetTimestamp() int64
	GetPushedReferences() PushedReferences
	Len() uint64
	ID(recompute ...bool) util.Bytes32
	BytesAndID(recompute ...bool) ([]byte, util.Bytes32)
	TxSize() uint
	BillableSize() uint64
	GetSize() uint64
	GetFee() util.String
}

// LocalRepo represents a local git repository on disk
type LocalRepo interface {
	LiteGit

	// GetName returns the name of the repo
	GetName() string

	// GetNameFromPath returns the name of the repo
	GetNameFromPath() string

	// GetNamespaceName returns the namespace this repo is associated to.
	GetNamespaceName() string

	// GetNamespace returns the repos's namespace
	GetNamespace() *state.Namespace

	// References returns an unsorted ReferenceIter for all references.
	References() (storer.ReferenceIter, error)

	// IsContributor checks whether a push key is a contributor to either
	// the repository or its namespace
	IsContributor(pushKeyID string) bool

	// GetRemoteURLs returns the remote URLS of the repository
	GetRemoteURLs() (urls []string)

	// DeleteObject deletes an object from a repository.
	DeleteObject(hash plumbing.Hash) error

	// Reference returns the reference for a given reference name.
	Reference(name plumbing.ReferenceName, resolved bool) (*plumbing.Reference, error)

	// Object returns an Object with the given hash.
	Object(t plumbing.ObjectType, h plumbing.Hash) (object.Object, error)

	// Objects returns an unsorted ObjectIter with all the objects in the repository.
	Objects() (*object.ObjectIter, error)

	// CommitObjects returns an unsorted ObjectIter with all the objects in the repository.
	CommitObjects() (object.CommitIter, error)

	// CommitObject returns a commit.
	CommitObject(h plumbing.Hash) (*object.Commit, error)

	// WrappedCommitObject returns commit that implements types.Commit interface.
	WrappedCommitObject(h plumbing.Hash) (Commit, error)

	// BlobObject returns a Blob with the given hash.
	BlobObject(h plumbing.Hash) (*object.Blob, error)

	// TagObject returns a Tag with the given hash.
	TagObject(h plumbing.Hash) (*object.Tag, error)

	// Tag returns a tag from the repository.
	Tag(name string) (*plumbing.Reference, error)

	// Config return the repository config
	Config() (*config.Config, error)

	// SetConfig sets the repo config
	SetConfig(cfg *config.Config) error

	// SetPath sets the repository root path
	SetPath(path string)

	// GetReferences returns all references in the repo
	GetReferences() (refs []plumbing.ReferenceName, err error)

	// GetPath returns the repository's path
	GetPath() string

	// GetState returns the repository's network state
	GetState() *state.Repository

	// SetState sets the repository's network state
	SetState(s *state.Repository)

	// Head returns the reference where HEAD is pointing to.
	Head() (string, error)

	// ObjectExist checks whether an object exist in the target repository
	ObjectExist(objHash string) bool

	// GetObjectSize returns the size of an object
	GetObjectSize(objHash string) (int64, error)

	// GetObjectDiskSize returns the size of the object as it exist on the system
	GetObjectDiskSize(objHash string) (int64, error)

	// GetEncodedObject returns an object
	GetEncodedObject(objHash string) (plumbing.EncodedObject, error)

	// WriteObjectToFile writes an object to the repository's objects directory
	WriteObjectToFile(objectHash string, content []byte) error

	// GetObject returns an object
	GetObject(objHash string) (object.Object, error)

	// GetCompressedObject compressed version of an object
	GetCompressedObject(hash string) ([]byte, error)

	// GetHost returns the storage engine of the repository
	GetHost() storage.Storer

	// Prune prunes objects older than the given time
	Prune(olderThan time.Time) error

	// NumIssueBranches counts the number of issues branches
	NumIssueBranches() (count int, err error)

	// GetAncestors returns the ancestors of the given commit up til the ancestor matching the stop hash.
	// The stop hash ancestor is not included in the result.
	// Reverse reverses the result
	GetAncestors(commit *object.Commit, stopHash string, reverse bool) (ancestors []*object.Commit, err error)
}

type LiteGit interface {
	RefDelete(refname string) error
	RefUpdate(refname, commitHash string) error
	TagDelete(tagname string) error
	RefGet(refname string) (string, error)
	GetRecentCommitHash() (string, error)
	GetHEAD(short bool) (string, error)
	NumCommits(branch string, noMerges bool) (int, error)
	GetConfig(path string) string
	CreateSignedEmptyCommit(msg, signingKey string, env ...string) error
	CreateTagWithMsg(args []string, msg, signingKey string, env ...string) error
	ListTreeObjects(treename string, recursive bool, env ...string) (map[string]string, error)
	ListTreeObjectsSlice(treename string, recursive, showTrees bool, env ...string) ([]string, error)
	RemoveEntryFromNote(notename, objectHash string, env ...string) error
	AddEntryToNote(notename, objectHash, note string, env ...string) error
	CreateBlob(content string) (string, error)
	UpdateRecentCommitMsg(msg, signingKey string, env ...string) error
	IsAncestor(commitA string, commitB string, env ...string) error
	HasMergeCommits(reference string, env ...string) (bool, error)
	GetMergeCommits(reference string, env ...string) ([]string, error)
	CreateSingleFileCommit(filename, content, commitMsg, parent string) (string, error)
	Checkout(refname string, create, force bool) error
	GetRefRootCommit(ref string) (string, error)
	GetRefCommits(ref string, noMerges bool) ([]string, error)
	Var(name string) (string, error)
	ExpandShortHash(hash string) (string, error)
}

// PushedReference represents a reference that was pushed by git client
type PushedReference struct {
	util.SerializerHelper `json:"-" msgpack:"-" mapstructure:"-"`
	Name                  string               `json:"name" msgpack:"name,omitempty"`       // The full name of the reference
	OldHash               string               `json:"oldHash" msgpack:"oldHash,omitempty"` // The hash of the reference before the push
	NewHash               string               `json:"newHash" msgpack:"newHash,omitempty"` // The hash of the reference after the push
	Nonce                 uint64               `json:"nonce" msgpack:"nonce,omitempty"`     // The next repo nonce of the reference
	Objects               []string             `json:"objects" msgpack:"objects,omitempty"` // A list of objects pushed to the reference
	MergeProposalID       string               `json:"mergeID" msgpack:"mergeID,omitempty"` // The merge proposal ID the reference is complaint with.
	Fee                   util.String          `json:"fee" msgpack:"fee,omitempty"`         // The merge proposal ID the reference is complaint with.
	PushSig               []byte               `json:"pushSig" msgpack:"pushSig,omitempty"` // The signature of from the push request token
	Data                  *types.ReferenceData `json:"data" msgpack:"data,omitempty"`       // Contains updates to the reference data
}

// EncodeMsgpack implements msgpack.CustomEncoder
func (pr *PushedReference) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(
		pr.Name,
		pr.OldHash,
		pr.NewHash,
		pr.Nonce,
		pr.Objects,
		pr.MergeProposalID,
		pr.Fee,
		pr.Data,
		pr.PushSig)
}

// DecodeMsgpack implements msgpack.CustomDecoder
func (pr *PushedReference) DecodeMsgpack(dec *msgpack.Decoder) error {
	return pr.DecodeMulti(dec,
		&pr.Name,
		&pr.OldHash,
		&pr.NewHash,
		&pr.Nonce,
		&pr.Objects,
		&pr.MergeProposalID,
		&pr.Fee,
		&pr.Data,
		&pr.PushSig)
}

// IsDeletable checks whether the pushed reference can be deleted
func (pr *PushedReference) IsDeletable() bool {
	return pr.NewHash == plumbing.ZeroHash.String()
}

// PushedReferences represents a collection of pushed references
type PushedReferences []*PushedReference

// GetByName finds a pushed reference by name
func (pf *PushedReferences) GetByName(name string) *PushedReference {
	for _, r := range *pf {
		if r.Name == name {
			return r
		}
	}
	return nil
}

// Names returns the names of the references
func (pf *PushedReferences) Names() (names []string) {
	for _, r := range *pf {
		names = append(names, r.Name)
	}
	return
}

// Commit represents a Commit.
type Commit interface {

	// NumParents returns the number of parents in a commit.
	NumParents() int

	// Parent returns the ith parent of a commit.
	Parent(i int) (Commit, error)

	// IsParent checks whether the specified hash is a parent of the commit
	IsParent(hash string) (bool, Commit)

	// UnWrap returns the underlying commit object
	UnWrap() *object.Commit

	// GetCommitter returns the one performing the commit, might be different from Author
	GetCommitter() *object.Signature

	// GetAuthor returns the original author of the commit.
	GetAuthor() *object.Signature

	// GetTreeHash returns the hash of the root tree of the commit
	GetTreeHash() plumbing.Hash

	// GetHash returns the hash of the commit object
	GetHash() plumbing.Hash

	// GetTree returns the tree from the commit
	GetTree() (*object.Tree, error)

	// File returns the file with the specified "path" in the commit and a
	// nil error if the file exists.
	File(path string) (*object.File, error)
}
