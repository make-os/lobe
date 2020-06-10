package types

import (
	"time"

	"gitlab.com/makeos/mosdef/types/state"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/storage"
)

// LocalRepo represents a local git repository on disk
type LocalRepo interface {
	LiteGit

	// GetName returns the name of the repo
	GetName() string

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

	// Tags return all tag references in the repository.
	// If you want to check to see if the tag is an annotated tag, you can call
	// TagObject on the hash Reference
	Tags() (storer.ReferenceIter, error)

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

	// IsAncestor checks whether commitA is an ancestor to commitB.
	// It returns ErrNotAncestor when not an ancestor.
	IsAncestor(commitA string, commitB string) error

	// SetPath sets the repository root path
	SetPath(path string)

	// GetReferences returns all references in the repo
	GetReferences() (refs []plumbing.ReferenceName, err error)

	// GetPath returns the repository's path
	GetPath() string

	// IsClean checks whether the working directory has no un-tracked, staged or modified files
	IsClean() (bool, error)

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

	// ObjectsOfCommit returns a hashes of objects a commit is composed of.
	// This objects a the commit itself, its tree and the tree blobs.
	ObjectsOfCommit(hash string) ([]plumbing.Hash, error)

	// GetObject returns an object
	GetObject(objHash string) (object.Object, error)

	// GetStorer returns the storage engine of the repository
	GetStorer() storage.Storer

	// Prune prunes objects older than the given time
	Prune(olderThan time.Time) error

	// NumIssueBranches counts the number of issues branches
	NumIssueBranches() (count int, err error)

	// GetAncestors returns the ancestors of the given commit up til the ancestor matching the stop hash.
	// The stop hash ancestor is not included in the result.
	// Reverse reverses the result
	GetAncestors(commit *object.Commit, stopHash string, reverse bool) (ancestors []*object.Commit, err error)
}

type RefFetchArgs struct {
	Remote    string
	RemoteRef string
	LocalRef  string
	Force     bool
	Verbose   bool
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
	HasMergeCommits(reference string, env ...string) (bool, error)
	GetMergeCommits(reference string, env ...string) ([]string, error)
	CreateSingleFileCommit(filename, content, commitMsg, parent string) (string, error)
	Checkout(refname string, create, force bool) error
	GetRefRootCommit(ref string) (string, error)
	GetRefCommits(ref string, noMerges bool) ([]string, error)
	Var(name string) (string, error)
	ExpandShortHash(hash string) (string, error)
	RefFetch(args RefFetchArgs) error
}

// Commit represents a Commit.
type Commit interface {
	object.Object

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

	// Tree returns the tree from the commit
	Tree() (*object.Tree, error)

	// File returns the file with the specified "path" in the commit and a
	// nil error if the file exists.
	File(path string) (*object.File, error)
}
