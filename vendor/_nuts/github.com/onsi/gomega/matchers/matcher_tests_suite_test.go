package matchers_test

import (
	"testing"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/gomega"
)

type myStringer struct {
	a string
}

func (s *myStringer) String() string {
	return s.a
}

type StringAlias string

type myCustomType struct {
	s   string
	n   int
	f   float32
	arr []string
}

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomega")
}
