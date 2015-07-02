package gitprocess_test

import (
	gp "com.mooregreatsoftware/go-git-process/lib"
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

	var _ = Context("Simple repo commit", func() {
		It("should create a simple commit", func() {
			tree, err := gp.AddPaths(tempRepo)
			commit, err := tempRepo.CreateCommit("", gp.Signature{}, gp.Signature{}, "test msg", tree, nil)
			Ω(err).ShouldNot(HaveOccurred())
			head, err := tempRepo.Head()
			Ω(err).ShouldNot(HaveOccurred())
			headOid := head.Oid()
			Ω(headOid).Should(Equal(commit.Oid()))
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
