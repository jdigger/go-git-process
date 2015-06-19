package gitprocess

import (
	"fmt"

	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"
)

// CreateFeatureBranch takes the name of the new feature branch to create and
// the path of the repository. It will fetch the latest from the server
// (if not running local-only), then creates a new branch based off of the
// configured "integration branch" and checks the repository out to it.
// The new branch's upstream is set to the integration branch.
// If the current branch was "_parking_" then the _parking_ branch is deleted.
func CreateFeatureBranch(name string, repoPath string) (*git.Repository, error) {
	repo, err := CreateRepository(repoPath)
	if err != nil {
		return nil, err
	}

	if err = repo.Fetch(FetchOptions{RemoteName: "origin", Prune: true}); err != nil {
		return nil, err
	}
	fmt.Printf("%v %v", repo, err)
	return nil, nil
}
