package gitprocess

import (
	"fmt"

	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"
)

// ******************************************
//
// Types
//
// ******************************************

// Index represents the current set of changes in git
type Index interface {
	AddByPath(path ...string) error
	WriteTree() (Tree, error)
}

// IndexSource provides a pointer to the current working git index
type IndexSource interface {
	Index() *Index
}

// ******************************************
//
// gitIndexSource
// implements the IndexSource interface
//
// ******************************************

type getIndexSourceStruct struct {
	gitRepo *git.Repository
}

func NewIndexSource(repo WorkingRepositoryReader) IndexSource {
	repoStruct, ok := repo.(RepositoryStruct)
	if ok {
		gr := git.Repository(*repoStruct.gitRepo())
		return getIndexSourceStruct{gitRepo: &gr}
	}
	gitRepo, err := git.OpenRepository(repo.Path())
	if err != nil {
		panic(fmt.Errorf("Problem opening repository at \"%s\": %v", repo.Path(), err))
	}
	return getIndexSourceStruct{gitRepo: gitRepo}
}

func (indexSource getIndexSourceStruct) Index() *Index {
	return indexFromGitRepo(gitRepository(*indexSource.gitRepo))
}

// ******************************************
//
// gitIndex
// implements the Index interface
//
// ******************************************

type gitIndexStruct struct {
	gitRepo  *git.Repository
	gitIndex *git.Index
}

func (gitIndex *gitIndexStruct) ensureGitIndex() {
	gitIndex.gitIndex, _ = gitIndex.gitRepo.Index()
}

func (gitIndex gitIndexStruct) AddByPath(paths ...string) error {
	gitIndex.ensureGitIndex()
	for _, path := range paths {
		if err := gitIndex.gitIndex.AddByPath(path); err != nil {
			return err
		}
	}
	return nil
}

func (gitIndex gitIndexStruct) WriteTree() (Tree, error) {
	gitIndex.ensureGitIndex()
	gitOid, err := gitIndex.gitIndex.WriteTree()
	if err != nil {
		return nil, err
	}
	oid := NewOid(gitOid.String())
	var tree Tree
	tree = treeStruct{oid: oid}
	return tree, nil
}
