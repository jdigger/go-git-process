package gitprocess

import "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"

// ******************************************
//
// Types
//
// ******************************************

// Tree represents a filesystem structure in git
type Tree interface {
	Oid() *Oid
}

// ******************************************
//
// treeStruct
// implements the Tree interface
//
// ******************************************

type treeStruct struct {
	oid *Oid
}

func (treeStruct treeStruct) Oid() *Oid {
	return treeStruct.oid
}

// ******************************************
//
// Functions
//
// ******************************************

// AddPaths adds the given paths to the Index and returns a Tree including them
func AddPaths(repo Repository, paths ...string) (Tree, error) {
	idx := *repo.Index()
	for _, path := range paths {
		err := idx.AddByPath(path)
		if err != nil {
			return nil, err
		}
	}
	return idx.WriteTree()
}

// gitTree translates a gitprocess.Tree to a git.Tree
func gitTree(tree Tree, gitRepo *git.Repository) (*git.Tree, error) {
	if tree == nil {
		return nil, nil
	}
	treeOid := tree.Oid()
	gitTreeOid := git.Oid(*treeOid)
	gitTree, err := gitRepo.LookupTree(&gitTreeOid)
	if err != nil {
		return nil, err
	}
	return gitTree, nil
}
