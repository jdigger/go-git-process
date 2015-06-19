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

/*
Repository encapsulates the primary operational parts of a git repository.
*/
type Repository interface {
	// See Fetcher.Fetch()
	Fetch(fetchOptions FetchOptions) error

	// the filesystem path of the repository
	Path() string

	Checkout(checkoutOptions CheckoutOptions) (*Branch, error)

	Index() *Index

	LookupTree(treeID *Oid) (*Tree, error)

	CreateCommit(refname string, author, committer *Signature, message string, tree *Tree, parents ...*Commit) (*Commit, error)

	Head() (*Commit, error)
}

type gitRepository git.Repository

func (gitRepo gitRepository) Index() *Index {
	gr := git.Repository(gitRepo)
	var idx Index
	idx = gitIndexStruct{gitRepo: &gr}
	return &idx
}

func (gitRepo gitRepository) Tree(treeID *Oid) (*Tree, error) {
	if treeID == nil {
		return nil, fmt.Errorf("treeID == nil")
	}
	var tree Tree
	tree = treeStruct{oid: treeID}
	return &tree, nil
}

/*
CreateRepository takes a path to a git repository and returns a instance
to work on it.
*/
func CreateRepository(repoPath string) (Repository, error) {
	log.Info("Opening repository at \"%s\"", repoPath)
	gitRepo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}

	repoStruct := &RepositoryStruct{RepoPath: gitRepo.Workdir()}
	repoStruct.RemoteFactory = func(remoteName string) (Remote, error) {
		return gitRemoteStruct{Repo: repoStruct, name: remoteName, gitRepo: gitRepo}, nil
	}
	var repo Repository = repoStruct
	repoStruct.RemotesFactory = func() (Remotes, error) {
		return listRemotes(gitRepo, &repo, &repoStruct.RemoteFactory)
	}
	repoStruct.Fetcher = func(options FetchOptions) error {
		return Fetch(repoStruct, options, repoStruct.RemoteFactory, repoStruct.RemotesFactory)
	}
	repoStruct.TreeFactory = func(treeID *Oid) (*Tree, error) {
		return gitRepository(*gitRepo).Tree(treeID)
	}
	repoStruct.Committer = func(refname string, author, committer *Signature, message string, tree *Tree, parents ...*Commit) (*Commit, error) {
		return createCommit(gitRepo, refname, author, committer, message, tree, parents...)
	}

	return repo, nil
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
type TreeFactory func(*Oid) (*Tree, error)

// Committer defines how to create a commit and returns the Oid created
type Committer func( /*refname*/ string /*author*/, *Signature /*committer*/, *Signature /*message*/, string /*tree*/, *Tree /*parents*/, ...*Commit) (*Commit, error)

// Fetch updates the local repositry. See Fetcher
func (repoStruct RepositoryStruct) Fetch(options FetchOptions) error {
	return repoStruct.Fetcher(options)
}

// Path gives the path on the filesystem to the repository
func (repoStruct RepositoryStruct) Path() string {
	return repoStruct.RepoPath
}

// Checkout checks the working directory out to the branch given in the options, returning the current Branch
func (repoStruct RepositoryStruct) Checkout(checkoutOptions CheckoutOptions) (*Branch, error) {
	// TODO: Implement
	return nil, nil
}

func (repoStruct RepositoryStruct) Index() *Index {
	return repoStruct.gitRepo().Index()
}

func (repoStruct RepositoryStruct) LookupTree(treeID *Oid) (*Tree, error) {
	return repoStruct.TreeFactory(treeID)
}

// CreateCommit creates a commit
func (repoStruct RepositoryStruct) CreateCommit(refname string, author *Signature, committer *Signature, message string, tree *Tree, parents ...*Commit) (*Commit, error) {
	return repoStruct.Committer(refname, author, committer, message, tree, parents...)
}

func (repoStruct RepositoryStruct) Head() (*Commit, error) {
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
	oid := Oid(*goid)
	commit := Commit{Oid: &oid}
	return &commit, nil
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
