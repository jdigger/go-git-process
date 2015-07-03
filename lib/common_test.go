package gitprocess_test

import (
	"fmt"

	gp "com.mooregreatsoftware/go-git-process/lib"
	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo"
	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/op/go-logging"

	"io/ioutil"
	"os"
	"path"
)

var log = logging.MustGetLogger("gitprocess_test")

type tempRepoType interface {
	gp.IndexSource
	gp.RepositoryWriter
	gp.RepositoryReader
	gp.WorkingRepositoryWriter
	gp.WorkingRepositoryReader
}

var tempRepo tempRepoType

func CleanupTestRepo(r tempRepoType) {
	err := os.RemoveAll(r.Path())
	CheckFatal(err)
}

func CheckFatal(err error) {
	if err != nil {
		Fail(err.Error())
	}
}

func CreateTestRepo() tempRepoType {
	// figure out where we can create the test repo
	path, err := ioutil.TempDir("", "git2go")
	CheckFatal(err)

	_, err = git.InitRepository(path, false)
	CheckFatal(err)

	tmpfile := "README"
	err = ioutil.WriteFile(path+"/"+tmpfile, []byte("foo\n"), 0644)
	CheckFatal(err)

	repo, err := gp.CreateRepository(path)
	CheckFatal(err)
	trepo, ok := repo.(tempRepoType)
	if !ok {
		Fail(fmt.Sprintf("Could not convert %#v to `tempRepoType`", repo))
	}
	return trepo
}

func SeedTestRepo(repo tempRepoType) gp.Commit {
	sig := gp.Signature{}

	tree, err := gp.AddPaths(repo, "README")
	CheckFatal(err)

	message := "This is a commit\n"
	commit, err := repo.CreateCommit("HEAD", sig, sig, message, tree)
	CheckFatal(err)

	return commit
}

func PathInRepo(repo *git.Repository, name string) string {
	return path.Join(path.Dir(path.Dir(repo.Path())), name)
}
