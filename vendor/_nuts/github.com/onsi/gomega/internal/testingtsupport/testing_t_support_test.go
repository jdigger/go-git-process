package testingtsupport_test

import (
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/gomega"

	"testing"
)

func TestTestingT(t *testing.T) {
	RegisterTestingT(t)
	Ω(true).Should(BeTrue())
}
