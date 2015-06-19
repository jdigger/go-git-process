package main

import "fmt"
import (
	"com.mooregreatsoftware/go-git-process/lib"
)

func main() {
	fmt.Println("Hello, 世界")
	gitprocess.CreateFeatureBranch("fooble", ".")
}
