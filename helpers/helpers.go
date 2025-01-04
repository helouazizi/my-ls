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
func ExtractContent(folderPath string) ([]os.FileInfo, bool, error) {
	dir, err := os.Open(folderPath)
	if err != nil {
		return nil, false, err
	}
	defer dir.Close()
	content, err := dir.Readdir(-1)
	if err != nil {
		return nil, false, err
	}
	FileInfo, err := os.Stat(folderPath)
	if err != nil {
		return nil, false, fmt.Errorf("ls: cannot access '%s': No such file or directory", folderPath)
	}
	if FileInfo.IsDir() {
		return content, true, nil
	}
	return nil, false, nil
}
func PrintContent(flags []string, content []os.FileInfo) []string {
	result := []string{}
	for _, v := range content {
		result = append(result, v.Name())
		fmt.Println(v.Name())
	}
	return result
}
func R(folderPath string) {
	content, isdir, err := ExtractContent(folderPath)
	if err != nil {
		fmt.Println(err)
	}
	folders := PrintContent(nil, content)
	if isdir {
		Ls(nil, folders)
	}

}

func Ls(flags, foldersPath []string) {
	Recursive := true
	for _, folderPath := range foldersPath {
		content, isdir, err := ExtractContent(folderPath)
		if !isdir {
			// in this case is a file
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			fmt.Println(folderPath, ":")
		}

		PrintContent(flags, content)
		if isdir && Recursive {
			R(folderPath)
		}
	}

}
