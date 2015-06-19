package gitprocess

import "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"

// ******************************************
//
// Types
//
// ******************************************

// Index represents the current set of changes in git
type Index interface {
	AddByPath(path ...string) error
	WriteTree() (*Tree, error)
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

func (gitIndex gitIndexStruct) WriteTree() (*Tree, error) {
	gitIndex.ensureGitIndex()
	gitOid, err := gitIndex.gitIndex.WriteTree()
	if err != nil {
		return nil, err
	}
	oid := Oid(*gitOid)
	var tree Tree
	tree = treeStruct{oid: &oid}
	return &tree, nil
}
