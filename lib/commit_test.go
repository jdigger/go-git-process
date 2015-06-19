package gitprocess_test

import (
	"com.mooregreatsoftware/go-git-process/lib"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/gomega"
)

var _ = Describe("Commit", func() {
	BeforeEach(func() {
		tempRepo = CreateTestRepo()
	})

	AfterEach(func() {
		CleanupTestRepo(tempRepo)
	})

	var _ = Context("Has remotes", func() {
		It("should do stuff", func() {
			tree, err := gitprocess.AddPaths(tempRepo)
			commit, err := tempRepo.CreateCommit("", nil, nil, "test msg", tree, nil)
			Ω(err).ShouldNot(HaveOccurred())
			head, err := tempRepo.Head()
			Ω(err).ShouldNot(HaveOccurred())
			headOid := *head.Oid
			Ω(headOid).Should(Equal(*commit.Oid))
		})

		// It("should fail fetching a failing repository", func() {
		// 	remoteFactory := func(remoteName string) (gitprocess.Remote, error) {
		// 		return fetchFailStub{RemoteStub{name: remoteName}}, nil
		// 	}
		//
		// 	err := gitprocess.Fetch(tempRepo, gitprocess.FetchOptions{}, remoteFactory, remoteListFactory)
		// 	Ω(err).Should(HaveOccurred())
		// })
	})
})

// *******************************************
//
// Stubs
//
// *******************************************
