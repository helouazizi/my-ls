// test/main.go
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	flags, foldersPath, err := Scan()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(foldersPath) == 0 {
		foldersPath = append(foldersPath, ".")
	}

	Ls(flags, foldersPath)
}

// Scan parses command-line arguments and extracts flags and folder paths.
func Scan() ([]string, []string, error) {
	foldersPath := []string{}
	flags := []string{}
	input := os.Args[1:]

	for _, arg := range input {
		if strings.HasPrefix(arg, "-") {
			exist, err := checkFlag(flags, arg)
			if err != nil {
				return nil, nil, err
			}
			if !exist {
				flags = append(flags, arg)
			}
		} else {
			foldersPath = append(foldersPath, arg)
		}
	}
	return flags, foldersPath, nil
}

// checkFlag validates flags and ensures they are not duplicated.
func checkFlag(flags []string, flag string) (bool, error) {
	for _, f := range flags {
		if f == flag {
			return true, nil
		}
	}

	// Validate flag format (e.g., no more than two dashes)
	dashCount := 0
	for _, char := range flag {
		if char == '-' {
			dashCount++
		}
		if dashCount > 2 {
			return false, fmt.Errorf("my-ls: unrecognized option '%s'", flag)
		}
	}
	return false, nil
}

// ExtractContent retrieves the contents of a directory or file.
func ExtractContent(folderPath string) ([]os.FileInfo, bool, error) {
	fileInfo, err := os.Stat(folderPath)
	if err != nil {
		return nil, false, fmt.Errorf("my-ls: cannot access '%s': No such file or directory", folderPath)
	}

	if !fileInfo.IsDir() {
		return []os.FileInfo{fileInfo}, false, nil
	}

	dir, err := os.Open(folderPath)
	if err != nil {
		return nil, false, err
	}
	defer dir.Close()

	content, err := dir.Readdir(-1)
	if err != nil {
		return nil, false, err
	}

	return content, true, nil
}

// PrintContent prints the contents of a directory based on the provided flags.
func PrintContent(flags []string, content []os.FileInfo) {
	// Apply sorting and filtering based on flags
	if containsFlag(flags, "-t") {
		sort.Slice(content, func(i, j int) bool {
			return content[i].ModTime().After(content[j].ModTime())
		})
	}
	if containsFlag(flags, "-r") {
		sort.Slice(content, func(i, j int) bool {
			return content[i].Name() > content[j].Name()
		})
	}

	for _, file := range content {
		if containsFlag(flags, "-a") || !strings.HasPrefix(file.Name(), ".") {
			if containsFlag(flags, "-l") {
				fmt.Printf("%s\t%d\t%s\t%s\n", file.Mode(), file.Size(), file.ModTime().Format(time.RFC3339), file.Name())
			} else {
				fmt.Println(file.Name())
			}
		}
	}
}

// R recursively lists the contents of directories.
func R(folderPath string, flags []string) {
	content, isDir, err := ExtractContent(folderPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isDir {
		fmt.Printf("\n%s:\n", folderPath)
		PrintContent(flags, content)
		for _, file := range content {
			if file.IsDir() {
				R(fmt.Sprintf("%s/%s", folderPath, file.Name()), flags)
			}
		}
	}
}

// Ls handles the main logic of listing files and directories.
func Ls(flags, foldersPath []string) {
	for _, folderPath := range foldersPath {
		content, isDir, err := ExtractContent(folderPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if isDir && len(foldersPath) > 1 {
			fmt.Printf("\n%s:\n", folderPath)
		}

		PrintContent(flags, content)

		if isDir && containsFlag(flags, "-R") {
			R(folderPath, flags)
		}
	}
}

// containsFlag checks if a specific flag is present in the flags slice.
func containsFlag(flags []string, flag string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}