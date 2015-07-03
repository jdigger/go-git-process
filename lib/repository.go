package gitprocess

import (
	"errors"
	"fmt"

	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"
)

// ******************************************
//
// Types
//
// ******************************************

// RepositoryReader encapsulates reading from a git repository
type RepositoryReader interface {
	LookupTree(treeID Oid) (Tree, error)
}

// WorkingRepositoryReader encapsulates reading from a non-bare git repository
type WorkingRepositoryReader interface {
	// the filesystem path of the repository
	Path() string

	Head() (Commit, error)
}

/*
RepositoryWriter encapsulates modifying a git repository.
*/
type RepositoryWriter interface {
	CreateCommit(refname string, author Signature, committer Signature, message string, tree Tree, parents ...Commit) (Commit, error)
}

/*
WorkingRepositoryWriter encapsulates modifying a non-bare git repository.
*/
type WorkingRepositoryWriter interface {
	// See Fetcher.Fetch()
	Fetch(fetchOptions FetchOptions) error

	Checkout(checkoutOptions CheckoutOptions) (Branch, error)
}

/*
Repository encapsulates the primary operational parts of a git repository.
*/
type Repository interface {
	RepositoryReader
	RepositoryWriter
}

/*
WorkingRepository encapsulates the primary operational parts of a git repository.
*/
type WorkingRepository interface {
	Repository
	WorkingRepositoryReader
	WorkingRepositoryWriter
}

/*
CreateRepository takes a path to a git repository and returns a instance
to work on it.
*/
func CreateRepository(repoPath string) (WorkingRepository, error) {
	log.Info("Opening repository at \"%s\"", repoPath)
	gitRepo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}

	repoStruct := &RepositoryStruct{RepoPath: gitRepo.Workdir()}
	repoStruct.RemoteFactory = func(remoteName string) (Remote, error) {
		return gitRemoteStruct{Repo: repoStruct, name: remoteName, gitRepo: gitRepo}, nil
	}
	var repo WorkingRepository = repoStruct
	repoStruct.RemotesFactory = func() (Remotes, error) {
		return listRemotes(*gitRepo, repo, repoStruct.RemoteFactory)
	}
	repoStruct.Fetcher = func(options FetchOptions) error {
		return Fetch(repoStruct, options, repoStruct.RemoteFactory, repoStruct.RemotesFactory)
	}
	repoStruct.TreeFactory = func(treeID Oid) (Tree, error) {
		return treeFromGitRepo(gitRepository(*gitRepo), treeID)
	}
	repoStruct.Committer = func(refname string, author Signature, committer Signature, message string, tree Tree, parents ...Commit) (Commit, error) {
		return createCommit(gitRepo, refname, author, committer, message, tree, parents...)
	}

	return repo, nil
}

type gitRepository git.Repository

func indexFromGitRepo(gitRepo gitRepository) *Index {
	gr := git.Repository(gitRepo)
	var idx Index
	idx = gitIndexStruct{gitRepo: &gr}
	return &idx
}

func treeFromGitRepo(gitRepo gitRepository, treeID Oid) (Tree, error) {
	if treeID == nil {
		return nil, fmt.Errorf("treeID == nil")
	}
	var tree Tree
	tree = treeStruct{oid: treeID}
	return tree, nil
}

// ******************************************
//
// RepositoryStruct
// implements Repository with fields that allow swapping out functionality
//
// ******************************************

// RepositoryStruct holds the data defining how to work with a Repository
type RepositoryStruct struct {
	_gitRepo       *gitRepository
	RepoPath       string
	Fetcher        Fetcher
	RemoteFactory  RemoteFactory
	RemotesFactory RemotesFactory
	TreeFactory    TreeFactory
	Committer      Committer
}

// Fetcher defines how to retrieve the latest content of a remote repository
type Fetcher func(options FetchOptions) error

// IndexFactory defines how to create an Index based on the current state of the working directory
type IndexFactory func() (*Index, error)

// TreeFactory defines how to create a Tree based on the given Index
type TreeFactory func(Oid) (Tree, error)

// Committer defines how to create a commit and returns the Oid created
type Committer func( /*refname*/ string /*author*/, Signature /*committer*/, Signature /*message*/, string /*tree*/, Tree /*parents*/, ...Commit) (Commit, error)

// Fetch updates the local repositry. See Fetcher
func (repoStruct RepositoryStruct) Fetch(options FetchOptions) error {
	return repoStruct.Fetcher(options)
}

// Path gives the path on the filesystem to the repository
func (repoStruct RepositoryStruct) Path() string {
	return repoStruct.RepoPath
}

// Checkout checks the working directory out to the branch given in the options, returning the current Branch
func (repoStruct RepositoryStruct) Checkout(checkoutOptions CheckoutOptions) (Branch, error) {
	// TODO: Implement
	return nil, nil
}

func (repoStruct RepositoryStruct) Index() *Index {
	return indexFromGitRepo(*repoStruct.gitRepo())
}

func (repoStruct RepositoryStruct) LookupTree(treeID Oid) (Tree, error) {
	return repoStruct.TreeFactory(treeID)
}

// CreateCommit creates a commit
func (repoStruct RepositoryStruct) CreateCommit(refname string, author Signature, committer Signature, message string, tree Tree, parents ...Commit) (Commit, error) {
	return repoStruct.Committer(refname, author, committer, message, tree, parents...)
}

func (repoStruct RepositoryStruct) Head() (Commit, error) {
	gitRepo := git.Repository(*repoStruct.gitRepo())
	ref, err := gitRepo.LookupReference("HEAD")
	if err != nil {
		return nil, err
	}
	ref, err = ref.Resolve()
	if err != nil {
		return nil, err
	}
	goid := ref.Target()
	if goid == nil {
		return nil, errors.New("goid == null")
	}
	commit := CreateCommit(NewOid(goid.String()))
	return commit, nil
}

func (repoStruct RepositoryStruct) gitRepo() *gitRepository {
	if repoStruct._gitRepo == nil {
		repoPath := repoStruct.Path()
		log.Info("Opening repository at \"%s\"", repoPath)
		gitRepo, err := git.OpenRepository(repoPath)
		if err != nil {
			panic(fmt.Errorf("Problem opening repository at \"%s\": %v", repoPath, err))
		}
		gr := gitRepository(*gitRepo)
		repoStruct._gitRepo = &gr
	}
	return repoStruct._gitRepo
}
