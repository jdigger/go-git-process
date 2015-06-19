package leafnodes_test

import (
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo/internal/leafnodes"
	. "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/gomega"

	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo/internal/codelocation"
	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/onsi/ginkgo/types"
)

var _ = Describe("It Nodes", func() {
	It("should report the correct type, text, flag, and code location", func() {
		codeLocation := codelocation.New(0)
		it := NewItNode("my it node", func() {}, types.FlagTypeFocused, codeLocation, 0, nil, 3)
		Ω(it.Type()).Should(Equal(types.SpecComponentTypeIt))
		Ω(it.Flag()).Should(Equal(types.FlagTypeFocused))
		Ω(it.Text()).Should(Equal("my it node"))
		Ω(it.CodeLocation()).Should(Equal(codeLocation))
		Ω(it.Samples()).Should(Equal(1))
	})
})
