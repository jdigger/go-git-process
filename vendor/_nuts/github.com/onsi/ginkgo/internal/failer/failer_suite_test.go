package failer_test

import (
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/gomega"

	"testing"
)

func TestFailer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Failer Suite")
}
