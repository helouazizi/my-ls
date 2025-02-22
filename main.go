// main.go
package main

import (
	"fmt"
	"os"

	"myls/helpers"
)

func main() {
	opts, directories, err := helpers.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, dir := range directories {
		if err := helpers.ListDirectory(dir, opts); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
