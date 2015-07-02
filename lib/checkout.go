package gitprocess

import "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"

type Branch interface{}

type CheckoutOptions struct {
	BranchName string
	branch     *Branch
}

type simpleBranch struct {
	name      string
	gitBranch *git.Branch
}

// Branch gives an instance to the Branch currently being used for HEAD
func (co *CheckoutOptions) Branch() *Branch {
	if co.branch == nil {
		b := Branch(simpleBranch{name: co.BranchName})
		co.branch = &b
	}
	return co.branch
}

func Checkout(repo WorkingRepository, branchName string) error {
	gitRepo, err := git.OpenRepository(repo.Path())
	if err != nil {
		return err
	}

	branch, err := gitRepo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return err
	}

	branchRef, err := branch.Resolve()
	if err != nil {
		return err
	}

	headRef, err := gitRepo.Head()
	if err != nil {
		return err
	}

	headRef.SetTarget(branchRef.Target(), nil, "")

	opts := &git.CheckoutOpts{Strategy: git.CheckoutSafeCreate}
	err = gitRepo.CheckoutHead(opts)
	return err
}
