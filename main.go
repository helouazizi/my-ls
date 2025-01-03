// main.go
package main

import (
	"fmt"
	"myls/helpers"
	"os"
)

func main() {
	// lets check the arguments
	if len(os.Args) > 3 {
		fmt.Println("Usage: my-ls <option> <folder-name>")
		return
	}

	folderPath := "./" // default folder
	option := ""
	if len(os.Args) == 3 {
		option = os.Args[1]
		folderPath = os.Args[2]
	} else if len(os.Args) == 2 {
		option = os.Args[1]
		folderPath = "./" // default folder
	}
	helpers.Ls(option, folderPath)

}
