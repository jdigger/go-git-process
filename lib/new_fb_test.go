package gitprocess_test

import (
	"com.mooregreatsoftware/go-git-process/lib"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/gomega"
)

var _ = Describe("New Feature Branch", func() {
	BeforeEach(func() {
		tempRepo = CreateTestRepo()
	})

	AfterEach(func() {
		CleanupTestRepo(tempRepo)
	})

	Describe("Fooble", func() {
		It("should create a feature branch", func() {
			_, err := gitprocess.CreateFeatureBranch("yep", tempRepo.Path())
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
