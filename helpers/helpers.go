// helpers/helpers.go
package helpers

import (
	"fmt"
	"os"
	"strings"
)

func Scan() ([]string, []string, error) {
	foldersPath := []string{}
	flags := []string{}
	input := os.Args
	// lets extract folderpath and the plags if exist
	for i := 1; i < len(input); i++ {
		index := input[i]
		if strings.HasPrefix(index, "-") {
			exist, err := checkFlag(flags, index)
			if err != nil {
				return nil, nil, err
			}
			if exist {
				continue
			}
			flags = append(flags, index)
		} else {
			foldersPath = append(foldersPath, index)
		}
	}
	return flags, foldersPath, nil

}
func checkFlag(falgs []string, flag string) (bool, error) {
	i := 0
	// check if the flag allready exist
	// if exist no need to do anything it not ann error

	for _, oldFlag := range falgs {
		if flag == oldFlag {
			return true, nil
		}
	}
	// check if the flag is valid
	// the flag must have at most 2 "--"
	for _, v := range flag {
		if i > 1 {
			return false, fmt.Errorf("ls: unrecognized option '%s'\nTry 'ls --help' for more information. ", flag)
		}
		if v == '-' {
			i++
		}
		/*if i == 2 && len(flag) > 3 {
			return false, fmt.Errorf("ls: unrecognized option '%s'\nTry 'ls --help' for more information. ", flag)
		}*/

	}
	return false, nil
}
func ExtractContent(folderPath string) ([]os.FileInfo, error) {
	_, err := os.Stat(folderPath)
	if err != nil {
		return nil, fmt.Errorf("ls: cannot access '%s': No such file or directory", folderPath)
	}
	dir, err := os.Open(folderPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	content, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return content, nil
}
func PrintContent(flags []string, content []os.FileInfo) {
	
}

func Ls(flags, foldersPath []string) {
	for _, folderPath := range foldersPath {
		content, err := ExtractContent(folderPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		PrintContent(flags, content)
	}

}
