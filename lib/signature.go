package gitprocess

import (
	"fmt"

	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"
)

// ******************************************
//
// Types
//
// ******************************************

// Signature is the "who & when" of an event in git
type Signature git.Signature

func (sig Signature) String() string {
	return fmt.Sprintf("Signature{\"%s\", \"%s\", %s}", sig.Name, sig.Email, sig.When.Local())
}
