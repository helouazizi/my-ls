// main.go
package main

import (
	"fmt"
	"myls/helpers"
	"os"
)

func main() {
	opts, directory := helpers.ParseFlags(os.Args[1:])
	if err := helpers.ListDirectory(directory, opts); err != nil {
		fmt.Println("Error:", err)
	}

}
