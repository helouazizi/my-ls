// helpers/helpers.go
package helpers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Ls(flag, folderPath string) {
	dir, err := os.Open(folderPath)
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
	fmt.Println("total", len(folderPath))
	switch flag {
	case "-l":
		L(content, true)
	case "-a":
		L(content, false)
	}

}

func L(content []os.FileInfo, condition bool) {
	for _, fi := range content {
		if condition {
			if strings.HasPrefix(fi.Name(), ".") || strings.HasPrefix(fi.Name(), "..") {
				continue
			}
		}

		fmt.Printf("%s %s %s %s %s\n", fi.Mode().String(), fi.Name(), strconv.Itoa(int(fi.Size())), fi.ModTime(), fi.ModTime().String())
	}
}
