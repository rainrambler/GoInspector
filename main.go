// GoParser project main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	//ParseGoSrc("demo.go.txt")
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <path>\n", os.Args[0])
		return
	}
	scanSrcDir(os.Args[1])

	fmt.Printf("INFO: Total Go files: %d\n", totalFiles)
}
