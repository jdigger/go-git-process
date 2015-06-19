package gitprocess_test

import (
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/gomega"
)

var _ = Describe("Checkout", func() {
	BeforeEach(func() {
		tempRepo = CreateTestRepo()
	})

	AfterEach(func() {
		CleanupTestRepo(tempRepo)
	})

	var _ = Context("Has remotes", func() {
		It("should do stuff", func() {
			SeedTestRepo(tempRepo)
			// err := gitprocess.Fetch(tempRepo, gitprocess.FetchOptions{}, remoteFactory, remoteListFactory)
			var err error
			Ω(err).ShouldNot(HaveOccurred())
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
