// helpers/helpers.go
package helpers

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type FileInfo struct {
	Name    string
	Mode    fs.FileMode
	Size    int64
	ModTime time.Time
	IsDir   bool
	Owner   string
	Group   string
}

type Options struct {
	Long      bool
	Recursive bool
	All       bool
	Reverse   bool
	TimeSort  bool
}

func Setfiles(directory string, options Options) ([]FileInfo, error) {
	dir_entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var files []FileInfo
	for _, file := range dir_entries {
		// lets ignore the hiden files id -a not exist
		if !options.All && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		//lets exytact file info
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		//nkow lets get ths file stat using syscall
		status, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return nil, errors.New("failed to get file stats")
		}
		// nkow we can extract user and group inforamtions
		user_info, _ := user.Lookup(strconv.Itoa(int(status.Uid)))
		grp_info, _ := user.Lookup(strconv.Itoa(int(status.Gid)))

		// nkow lets fill the fileinfo struct we all information about this dir
		files = append(files, FileInfo{
			Name:    info.Name(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			Owner:   user_info.Username,
			Group:   grp_info.Username,
		})

	}
	SortFiles(&files, options)

	return files, nil
}
func SortFiles(files *[]FileInfo, opts Options) {
	if opts.TimeSort {
		sort.Slice(*files, func(i, j int) bool {
			return (*files)[i].ModTime.After((*files)[j].ModTime)
		})
	} else {
		sort.Slice(*files, func(i, j int) bool {
			return (*files)[i].Name < (*files)[j].Name
		})
	}

	if opts.Reverse {
		sort.SliceStable(*files, func(i, j int) bool {
			return !((*files)[i].ModTime.Before((*files)[j].ModTime))
		})
	}
}

// Scan parses command-line arguments and extracts flags and folder paths.
func Scan() (string, []string, error) {
	foldersPath := []string{}
	flags := ""
	input := os.Args[1:]

	for _, arg := range input {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			flags += string(arg)
		} else {
			foldersPath = append(foldersPath, arg)
		}
	}

	status, badflag := checkFlag(flags)
	if !status {
		return "nil", nil, fmt.Errorf("my-ls: invalid option -- '%s'\nTry 'ls --help' for more information", badflag)
	}

	return flags, foldersPath, nil
}

// checkFlag validates flags and ensures they are not duplicated.
func checkFlag(flags string) (bool, string) {
	patern := "alrRt-"

	for _, f := range flags {
		if !strings.Contains(patern, string(f)) {
			return false, string(f)
		}
	}

	return true, "test"
}

// containsFlag checks if a specific flag is present in the flags slice.
func containsFlag(flags string, flag string) bool {
	for _, f := range flags {
		if string(f) == flag {
			return true
		}
	}
	return false
}

// ExtractContent retrieves the contents of a directory or file.
func ExtractContent(folderPath string) ([]os.FileInfo, bool, error) {
	fileInfo, err := os.Stat(folderPath)
	if err != nil {
		return nil, false, fmt.Errorf("my-ls-1: cannot access '%s': No such file or directory", folderPath)
	}
	// tis condition for if it is a file
	if !fileInfo.IsDir() {
		return []os.FileInfo{fileInfo}, false, nil
	}

	dir, err := os.Open(folderPath)
	if err != nil {
		// fmt.Println("here")
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
func PrintContent(flagsGrouped string, content []os.FileInfo) {
	// Apply sorting and filtering based on flags
	// lets sort the content by time if -t exist
	if containsFlag(flagsGrouped, "t") {
		sort.Slice(content, func(i, j int) bool {
			return content[i].ModTime().After(content[j].ModTime())
		})
	}
	if containsFlag(flagsGrouped, "r") {
		sort.Slice(content, func(i, j int) bool {
			return content[i].Name() > content[j].Name()
		})
	}

	if containsFlag(flagsGrouped, "l") {
		// lets print the total of the folder size
		total := 0
		for _, file := range content {
			total += int(file.Size() / 512)
		}
		fmt.Println("total: ", total)
	}

	for _, file := range content {
		if containsFlag(flagsGrouped, "l") {
			if containsFlag(flagsGrouped, "a") {
				fmt.Printf("%s  %d  %s  %s\n", file.Mode(), file.Size(), file.ModTime().Format(time.RFC3339), file.Name())
			} else {
				if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "..") {
					continue
				}
				fmt.Printf("%s  %d  %s  %s\n", file.Mode(), file.Size(), file.ModTime().Format(time.RFC3339), file.Name())
			}
		} else {
			if containsFlag(flagsGrouped, "a") {
				fmt.Print(file.Name(), "  ")
			} else {
				if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "..") {
					continue
				}
				fmt.Print(file.Name(), "  ")
			}
		}
	}
	fmt.Println()
}

// R recursively lists the contents of directories.
func R(folderPath string, flagsGrouped string, i int) {
	test := strings.Split(folderPath, "/")
	lastone := test[len(test)-1]
	if !containsFlag(flagsGrouped, "a") && strings.HasPrefix(lastone, ".") && lastone != "." && lastone != ".." {
		return
	}
	content, isDir, err := ExtractContent(folderPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isDir {
		if i == 0 {
			fmt.Printf("%s:\n", folderPath)
		} else {
			fmt.Printf("\n%s:\n", folderPath)
		}

		PrintContent(flagsGrouped, content)
		for _, file := range content {
			if file.IsDir() {
				R(fmt.Sprintf("%s/%s", folderPath, file.Name()), flagsGrouped, -1111)
			}
		}
	} else {
		PrintContent(flagsGrouped, content)
	}
}

// Ls handles the main logic of listing files and directories.
func Ls(flags string, foldersPath []string) {
	i := 0
	for _, folderPath := range foldersPath {
		content, _, err := ExtractContent(folderPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if containsFlag(flags, "R") {
			R(folderPath, flags, i)
			i++
		} else {
			PrintContent(flags, content)
		}
	}
}
