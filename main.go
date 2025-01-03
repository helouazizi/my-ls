// main.go
package main

import (
	"fmt"
	"myls/helpers"
)

func main() {
	falgs, foldersPath, err := helpers.Scan()
	if err != nil {
		fmt.Println(err)
		return
	}
	helpers.Ls(falgs, foldersPath)

}
