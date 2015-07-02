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

// Oid is an Object ID in git
type Oid interface {
	fmt.Stringer
}

// NewOid creates a new Oid for the given string. It panics if it's not a valid oid.
func NewOid(oidStr string) Oid {
	err := validateOid(oidStr)
	if err != nil {
		panic(err)
	}

	return oidStruct{oidStr: oidStr}
}

// ******************************************
//
// oidStruct
// implements the Oid interface
//
// ******************************************

type oidStruct struct {
	oidStr string
}

func (oid oidStruct) String() string {
	return oid.oidStr
}

// ******************************************
//
// Functions
//
// ******************************************

func toGitOid(oid Oid) (*git.Oid, error) {
	return git.NewOid(oid.String())
}

func validateOid(oidStr string) error {
	_, err := git.NewOid(oidStr)
	return err
}
