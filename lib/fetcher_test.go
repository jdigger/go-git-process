package gitprocess_test

import (
	"fmt"

	"com.mooregreatsoftware/go-git-process/lib"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/gomega"
)

var _ = Describe("Fetch", func() {
	BeforeEach(func() {
		tempRepo = CreateTestRepo()
	})

	AfterEach(func() {
		CleanupTestRepo(tempRepo)
	})

	remoteFactory := func(remoteName string) (gitprocess.Remote, error) {
		return RemoteStub{name: remoteName}, nil
	}

	var _ = Context("Has remotes", func() {
		remoteListFactory := func() (gitprocess.Remotes, error) {
			return gitprocess.Remotes{RemoteStub{name: "origin"}}, nil
		}

		It("should fetch a remotable repository", func() {
			err := gitprocess.Fetch(tempRepo, gitprocess.FetchOptions{}, remoteFactory, remoteListFactory)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should fail fetching a failing repository", func() {
			remoteFactory := func(remoteName string) (gitprocess.Remote, error) {
				return fetchFailStub{RemoteStub{name: remoteName}}, nil
			}

			err := gitprocess.Fetch(tempRepo, gitprocess.FetchOptions{}, remoteFactory, remoteListFactory)
			Ω(err).Should(HaveOccurred())
		})
	})

	var _ = Context("Has no remotes", func() {
		remoteListFactory := func() (gitprocess.Remotes, error) {
			return nil, nil
		}

		It("should be a no-op doing a fetch", func() {
			err := gitprocess.Fetch(tempRepo, gitprocess.FetchOptions{}, remoteFactory, remoteListFactory)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})

// *******************************************
//
// Stubs
//
// *******************************************

type RemoteStub struct {
	name string
}

func (remote RemoteStub) Prune() error {
	return nil
}

func (remote RemoteStub) Name() string {
	return remote.name
}

func (remote RemoteStub) Fetch() error {
	return nil
}

type fetchFailStub struct {
	RemoteStub
}

func (remote fetchFailStub) Fetch() error {
	log.Error("Going BOOM")
	return fmt.Errorf("Boom")
}
