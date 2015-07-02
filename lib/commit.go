package gitprocess

import (
	"errors"
	"fmt"
	"time"

	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"
)

// ******************************************
//
// Types
//
// ******************************************

// Commit is a git commit
type Commit struct {
	Oid *Oid
}

func (commit Commit) String() string {
	return "Commit{" + commit.Oid.String() + "}"
}

// ******************************************
//
// Functions
//
// ******************************************

// DefaultSignature returns a signature for the repository's user and "now"
func defaultSignature(gitRepo *git.Repository) Signature {
	config, err := gitRepo.Config()
	if err != nil {
		panic(err)
	}

	// forEachConfigEntry(config, func(configEntry git.ConfigEntry) {
	// 	log.Info("%#v", configEntry)
	// })

	userName, err := config.LookupString("user.name")
	if err != nil {
		panic(errors.New("Could not find a configuration value for \"user.name\""))
	}

	userEmail, err := config.LookupString("user.email")
	if err != nil {
		panic(errors.New("Could not find a configuration value for \"user.email\""))
	}

	return Signature{
		Name:  userName,
		Email: userEmail,
		When:  time.Now(),
	}
}

func forEachConfigEntry(config *git.Config, processor func(configEntry git.ConfigEntry)) {
	configIterator, err := config.NewIterator()
	if err != nil {
		panic(err)
	}
	configEntry, err := configIterator.Next()
	for err == nil {
		processor(*configEntry)
		configEntry, err = configIterator.Next()
	}
}

func createCommit(gitRepo *git.Repository, refname string, author Signature, committer Signature, message string, tree Tree, parents ...Commit) (Commit, error) {
	if committer.Email == "" {
		sig := defaultSignature(gitRepo)
		committer = sig
	}

	if author.Email == "" {
		author = committer
	}

	if refname == "" {
		// The branch to update with the new commit.
		// If it's a symbolic reference - like HEAD typically is - it updates what it points to.
		// If the reference (branch) does not yet exist, it is created.
		refname = "HEAD"
	}

	if tree == nil {
		return Commit{}, errors.New("tree == nil")
	}

	log.Info("Creating commit on %s with %s and %s", refname, committer, author)

	gitAuthorSig := git.Signature(author)
	gitCommitterSig := git.Signature(committer)

	gitTree, err := gitTree(tree, gitRepo)
	if err != nil {
		return Commit{}, err
	}

	gitParents := [](*git.Commit){}
	for _, parent := range parents {
		if parent.Oid != nil {
			gitParents = append(gitParents, toGitCommit(gitRepo, parent))
		}
	}

	gitOid, err := gitRepo.CreateCommit(refname, &gitAuthorSig, &gitCommitterSig, message, gitTree, gitParents...)
	if err != nil {
		return Commit{}, err
	}
	oid := Oid(*gitOid)

	return Commit{Oid: &oid}, nil
}

func toGitCommit(gitRepo *git.Repository, commit Commit) *git.Commit {
	oid := *commit.Oid
	goid := git.Oid(oid)
	gitCommit, err := gitRepo.LookupCommit(&goid)
	if err != nil {
		panic(err)
	}
	if gitCommit == nil {
		panic(fmt.Errorf("Could not find commit %s in %s", goid, gitRepo.Workdir()))
	}
	return gitCommit
}
