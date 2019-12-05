package repo

import (
	"crypto/rsa"
	"fmt"
	"gopkg.in/src-d/go-git.v4/config"
	"os"
	"strings"

	"github.com/makeos/mosdef/crypto"
	"github.com/makeos/mosdef/types"
	"github.com/makeos/mosdef/util"
	"github.com/thoas/go-funk"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"

	"github.com/pkg/errors"
)

// Repo represents a git repository
type Repo struct {
	git   *git.Repository
	db    DBOperations
	path  string
	name  string
	ops   *GitOps
	state *types.Repository
}

// State returns the repository's network state
func (r *Repo) State() *types.Repository {
	return r.state
}

// Path returns the repository's path
func (r *Repo) Path() string {
	return r.path
}

// References returns an unsorted ReferenceIter for all references.
func (r *Repo) References() (storer.ReferenceIter, error) {
	return r.git.References()
}

// RefDelete executes `git update-ref -d <refname>` to delete a reference
func (r *Repo) RefDelete(refname string) error {
	return r.ops.RefDelete(refname)
}

// RefUpdate executes `git update-ref <refname> <commit hash>` to update/create a reference
func (r *Repo) RefUpdate(refname, commitHash string) error {
	return r.ops.RefUpdate(refname, commitHash)
}

// RefGet returns the hash content of a reference.
func (r *Repo) RefGet(refname string) (string, error) {
	return r.ops.RefGet(refname)
}

// TagDelete executes `git tag -d <tagname>` to delete a tag
func (r *Repo) TagDelete(tagname string) error {
	return r.ops.TagDelete(tagname)
}

// ListTreeObjects executes `git tag -d <tagname>` to delete a tag
func (r *Repo) ListTreeObjects(treename string, recursive bool, env ...string) (map[string]string,
	error) {
	return r.ops.ListTreeObjects(treename, recursive, env...)
}

// DeleteObject deletes an object from a repository.
func (r *Repo) DeleteObject(hash plumbing.Hash) error {
	return r.git.DeleteObject(hash)
}

// Reference deletes an object from a repository.
func (r *Repo) Reference(name plumbing.ReferenceName, resolved bool) (*plumbing.Reference, error) {
	return r.git.Reference(name, resolved)
}

// Object returns an Object with the given hash.
func (r *Repo) Object(t plumbing.ObjectType, h plumbing.Hash) (object.Object, error) {
	return r.git.Object(t, h)
}

// Objects returns an unsorted ObjectIter with all the objects in the repository.
func (r *Repo) Objects() (*object.ObjectIter, error) {
	return r.git.Objects()
}

// CommitObjects returns an unsorted ObjectIter with all the objects in the repository.
func (r *Repo) CommitObjects() (object.CommitIter, error) {
	return r.git.CommitObjects()
}

// CommitObject returns an unsorted ObjectIter with all the objects in the repository.
func (r *Repo) CommitObject(h plumbing.Hash) (*object.Commit, error) {
	return r.git.CommitObject(h)
}

// BlobObject returns a Blob with the given hash.
func (r *Repo) BlobObject(h plumbing.Hash) (*object.Blob, error) {
	return r.git.BlobObject(h)
}

// TagObject returns a Tag with the given hash.
func (r *Repo) TagObject(h plumbing.Hash) (*object.Tag, error) {
	return r.git.TagObject(h)
}

// Tag returns a tag from the repository.
func (r *Repo) Tag(name string) (*plumbing.Reference, error) {
	return r.git.Tag(name)
}

// Config return the repository config
func (r *Repo) Config() (*config.Config, error) {
	return r.git.Config()
}

// GetConfig finds and returns a config value
func (r *Repo) GetConfig(path string) string {
	return r.ops.GetConfig(path)
}

// GetRecentCommit gets the hash of the recent commit.
// Returns ErrNoCommits if no commits exist
func (r *Repo) GetRecentCommit() (string, error) {
	return r.ops.GetRecentCommit()
}

// UpdateRecentCommitMsg updates the recent commit message.
// msg: The commit message which is passed to the command's stdin.
// signingKey: The signing key
// env: Optional environment variables to pass to the command.
func (r *Repo) UpdateRecentCommitMsg(msg, signingKey string, env ...string) error {
	return r.ops.UpdateRecentCommitMsg(msg, signingKey, env...)
}

// CreateTagWithMsg an annotated tag.
// args: `git tag` options (NOTE: -a and --file=- are added by default)
// msg: The tag's message which is passed to the command's stdin.
// signingKey: The signing key to use
// env: Optional environment variables to pass to the command.
func (r *Repo) CreateTagWithMsg(args []string, msg, signingKey string, env ...string) error {
	return r.ops.CreateTagWithMsg(args, msg, signingKey, env...)
}

// RemoveEntryFromNote removes a note
func (r *Repo) RemoveEntryFromNote(notename, objectHash string, env ...string) error {
	return r.ops.RemoveEntryFromNote(notename, objectHash, env...)
}

// CreateBlob creates a blob object
func (r *Repo) CreateBlob(content string) (string, error) {
	return r.ops.CreateBlob(content)
}

// AddEntryToNote adds a note
func (r *Repo) AddEntryToNote(notename, objectHash, note string, env ...string) error {
	return r.AddEntryToNote(notename, objectHash, note, env...)
}

// ListTreeObjectsSlice returns a slice containing objects name of tree entries
func (r *Repo) ListTreeObjectsSlice(treename string, recursive, showTrees bool,
	env ...string) ([]string, error) {
	return r.ops.ListTreeObjectsSlice(treename, recursive, showTrees, env...)
}

// GetName returns the name of the repo
func (r *Repo) GetName() string {
	return r.name
}

func getRepo(path string) (types.BareRepo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repo{
		git:  repo,
		path: path,
	}, nil
}

func getRepoWithGitOpt(gitBinPath, path string) (types.BareRepo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repo{
		ops:  NewGitOps(gitBinPath, path),
		git:  repo,
		path: path,
	}, nil
}

// getCurrentWDRepo returns a Repo instance pointed to the repository
// in the current working directory.
func getCurrentWDRepo(gitBinDir string) (types.BareRepo, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current working directory")
	}

	// Since we expect the working directory to be a git working tree,
	// we need to get a repo instance to verify it
	repo, err := getRepoWithGitOpt(gitBinDir, wd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open repository")
	} else if repoCfg, _ := repo.Config(); repoCfg.Core.IsBare {
		return nil, errors.New("expected a working tree. this is a bare repository")
	}

	return repo, nil
}

// AmendRecentCommitTxLineCmd attempts to add or amend transaction argument to
// the recent commit. If no transaction line exist, it will add a new
// one populated with the provided arguments.
// Transaction line allows the user to set transaction arguments such as
// fee, public key, etc. A txline has the format tx: fee=10, pk=ad1..xyz, nonce=1
func AmendRecentCommitTxLineCmd(gitBinDir, txFee, txNonce, signingKey string) error {

	repo, err := getCurrentWDRepo(gitBinDir)
	if err != nil {
		return err
	}

	// Get the signing key id from the git config if not passed as an argument
	if signingKey == "" {
		signingKey = repo.GetConfig("user.signingKey")
	}
	if signingKey == "" {
		return errors.New("signing key was not set or provided")
	}

	// Get recent commit hash of the current branch
	hash, err := repo.GetRecentCommit()
	if err != nil {
		if err == ErrNoCommits {
			return errors.New("no commits have been created yet")
		}
		return err
	}

	commit, _ := repo.CommitObject(plumbing.NewHash(hash))
	msg := util.RemoveTxLine(commit.Message)

	// Get the public key of the signing key
	pkEntity, err := crypto.GetGPGPublicKey(signingKey, repo.GetConfig("gpg.program"), "")
	if err != nil {
		return errors.Wrap(err, "failed to get gpg public key")
	}

	// Get the public key network ID
	pkID := util.RSAPubKeyID(pkEntity.PrimaryKey.PublicKey.(*rsa.PublicKey))

	// TODO:If nonce is not provided, get nonce from a --source (remote node or local
	// node). Requires an RPC endpoint

	// Construct the tx line and append to the current message
	txLine := util.MakeTxLine(txFee, txNonce, pkID, nil)
	msg += "\n\n" + txLine

	// Update the recent commit message
	if err = repo.UpdateRecentCommitMsg(msg, signingKey); err != nil {
		return err
	}

	return nil
}

// CreateTagWithTxLineCmd creates a tag and adds a txline to the tag message
func CreateTagWithTxLineCmd(args []string, gitBinDir, txFee, txNonce, signingKey string) error {

	repo, err := getCurrentWDRepo(gitBinDir)
	if err != nil {
		return err
	}

	parsed := util.ParseSimpleArgs(args)

	// If -u flag is provided in the git args, use it a signing key
	if parsed["u"] != "" {
		signingKey = parsed["u"]
	}
	// Get the signing key id from the git config if not passed via app -u flag
	if signingKey == "" {
		signingKey = repo.GetConfig("user.signingKey")
	}
	// Return error if we still don't have a signing key
	if signingKey == "" {
		return errors.New("signing key was not set or provided")
	}

	// Get the user-supplied message from the arguments provided
	msg := ""
	if m, ok := parsed["m"]; ok {
		msg = m
	} else if message, ok := parsed["message"]; ok {
		msg = message
	}

	// Remove -m or --message flag from args
	args = util.RemoveFlagVal(args, []string{"m", "message", "u"})

	// Get the public key of the signing key
	pkEntity, err := crypto.GetGPGPublicKey(signingKey, repo.GetConfig("gpg.program"), "")
	if err != nil {
		return errors.Wrap(err, "failed to get gpg public key")
	}

	// Get the public key network ID
	pkID := util.RSAPubKeyID(pkEntity.PrimaryKey.PublicKey.(*rsa.PublicKey))

	// TODO:If nonce is not provided, get nonce from a --source (remote node or local
	// node). Requires an RPC endpoint

	// Construct the tx line and append to the current message
	txLine := util.MakeTxLine(txFee, txNonce, pkID, nil)
	msg += "\n\n" + txLine

	// Create the tag
	if err = repo.CreateTagWithMsg(args, msg, signingKey); err != nil {
		return err
	}

	return nil
}

// AddSignedTxBlobCmd creates a blob object that contains a signed tx line.
func AddSignedTxBlobCmd(gitBinDir, txFee, txNonce, signingKey, note string) error {

	repo, err := getCurrentWDRepo(gitBinDir)
	if err != nil {
		return err
	}

	// Get the signing key id from the git config if not provided via -s flag
	if signingKey == "" {
		signingKey = repo.GetConfig("user.signingKey")
	}
	// Return error if we still don't have a signing key
	if signingKey == "" {
		return errors.New("signing key was not set or provided")
	}

	// Enforce the inclusion of `refs/notes` to the note argument
	if !strings.HasPrefix("refs/notes", note) {
		note = "refs/notes/" + note
	}

	// Find a list of all notes entries in the note
	noteEntries, err := repo.ListTreeObjects(note, false)
	if err != nil {
		msg := fmt.Sprintf("unable to fetch note entries for tree object (%s)", note)
		return errors.Wrap(err, msg)
	}

	// From the entries, find existing tx blob and stop after the first one
	var lastTxBlob *object.Blob
	for hash := range noteEntries {
		obj, err := repo.BlobObject(plumbing.NewHash(hash))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to read object (%s)", hash))
		}
		r, err := obj.Reader()
		if err != nil {
			return err
		}
		prefix := make([]byte, 3)
		r.Read(prefix)
		if string(prefix) == util.TxLinePrefix {
			lastTxBlob = obj
			break
		}
	}

	// Remove the last tx blob from the note, if present
	if lastTxBlob != nil {
		err = repo.RemoveEntryFromNote(note, noteEntries[lastTxBlob.Hash.String()])
		if err != nil {
			return errors.Wrap(err, "failed to delete existing transaction blob")
		}
	}

	// Get the commit hash the note is currently referencing.
	// We need to add this hash to the signature.
	noteRef, err := repo.Reference(plumbing.ReferenceName(note), true)
	if err != nil {
		return errors.Wrap(err, "failed to get note reference")
	}
	noteHash := noteRef.Hash().String()

	// Get the public key of the signing key
	pkEntity, err := crypto.GetGPGPrivateKey(signingKey, repo.GetConfig("gpg.program"), "")
	if err != nil {
		return errors.Wrap(err, "failed to get gpg public key")
	}

	// Get the public key network ID
	pkID := util.RSAPubKeyID(pkEntity.PrimaryKey.PublicKey.(*rsa.PublicKey))

	// TODO:If nonce is not provided, get nonce from a --source (remote node or local
	// node). Requires an RPC endpoint

	// Sign a message composed of the tx information
	// fee + nonce + public key id + note hash
	sigMsg := []byte(txFee + txNonce + pkID + noteHash)
	sig, err := crypto.GPGSign(pkEntity, sigMsg)
	if err != nil {
		return errors.Wrap(err, "failed to sign transaction parameters")
	}

	// Construct the tx line
	txLine := util.MakeTxLine(txFee, txNonce, pkID, sig)

	// Create a blob with 0 byte content which be the subject of our note.
	blobHash, err := repo.CreateBlob("")
	if err != nil {
		return err
	}

	// Next we add the tx blob to the note
	if err = repo.AddEntryToNote(note, blobHash, txLine); err != nil {
		return errors.Wrap(err, "failed to add tx blob")
	}

	fmt.Println("Added signed transaction object to note")

	return nil
}

// getTreeEntries returns all entries in a tree.
func getTreeEntries(repo types.BareRepo, treeHash string) ([]string, error) {
	entries, err := repo.ListTreeObjectsSlice(treeHash, true, true)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// getCommitHistory gets all objects that led up to the given commit, such
// as parent commits, trees and blobs.
// repo: The target repository
// commit: The target commit
// stopCommitHash: A commit hash that when found triggers the end of the search.
func getCommitHistory(repo types.BareRepo, commit *object.Commit, stopCommitHash string) ([]string, error) {
	var hashes []string

	// Stop if commit hash matches the stop hash
	if commit.Hash.String() == stopCommitHash {
		return hashes, nil
	}

	// Add the commit and the tree hash
	hashes = append(hashes, commit.Hash.String())
	hashes = append(hashes, commit.TreeHash.String())

	// Get entries of the tree (blobs and sub-trees)
	entries, err := getTreeEntries(repo, commit.TreeHash.String())
	if err != nil {
		return nil, err
	}
	hashes = append(hashes, entries...)

	// Perform same operation on the parents of the commit
	err = commit.Parents().ForEach(func(parent *object.Commit) error {
		childHashes, err := getCommitHistory(repo, parent, stopCommitHash)
		if err != nil {
			return err
		}
		hashes = append(hashes, childHashes...)
		return nil
	})

	return funk.UniqString(hashes), err
}

// getObjectsSize returns the total size of the given objects in the repo.
// Panics if unable to find an object.
func getObjectsSize(repo types.BareRepo, objects []string) uint64 {
	var size int64
	for _, hash := range objects {
		obj, err := repo.Object(plumbing.AnyObject, plumbing.NewHash(hash))
		if err != nil {
			panic(err)
		}
		encoded := &plumbing.MemoryObject{}
		if err = obj.Encode(encoded); err != nil {
			panic(err)
		}
		size += encoded.Size()
	}
	return uint64(size)
}