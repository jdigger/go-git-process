package gitprocess

import "com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"

// ******************************************
//
// Types
//
// ******************************************

// Oid is an Object ID in git
type Oid git.Oid

func (oid Oid) String() string {
	goid := git.Oid(oid)
	return goid.String()
}

// Equal compares two Oids for equality
func (oid Oid) Equal(other Oid) bool {
	goid := git.Oid(oid)
	gother := git.Oid(other)
	return goid.Equal(&gother)
}
