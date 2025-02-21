// main.go
package main

import (
	"fmt"
	"myls/helpers"
	"os"
)

func main() {
	opts, directories := helpers.ParseFlags(os.Args[1:])
	for _, dir := range directories {
		if err := helpers.ListDirectory(dir, opts); err != nil {
			fmt.Println("Error:", err)
		}

	}

}
