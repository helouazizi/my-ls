// main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	// lets check the arguments
	if len(os.Args) > 3 {
		fmt.Println("Usage: my-ls <option> <folder-name>")
		return
	}

	folder := "./"
	if len(os.Args) == 3 {
		folder = os.Args[2]
	}
	dir, err := os.Open(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dir.Close()
	content, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("total", len(folder))
	for _, fi := range content {
		fmt.Println(fi.Name())
	}
}
