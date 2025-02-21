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

func ParseFlags(args []string) (Options, string) {
	opts := Options{}
	dir := "."
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, char := range arg[1:] {
				switch char {
				case 'l':
					opts.Long = true
				case 'R':
					opts.Recursive = true
				case 'a':
					opts.All = true
				case 'r':
					opts.Reverse = true
				case 't':
					opts.TimeSort = true
				}
			}
		} else {
			dir = arg
		}
	}
	return opts, dir
}

func Getfiles(directory string, options Options) ([]FileInfo, error) {
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
func PrintFiles(files []FileInfo, opts Options) {
	for _, file := range files {
		if opts.Long {
			fmt.Printf("%s %s %s %10d %s %s\n", file.Mode, file.Owner, file.Group, file.Size, file.ModTime.Format("Jan 02 15:04"), file.Name)
		} else {
			fmt.Println(file.Name)
		}
	}
}
func ListDirectory(directory string, opts Options) error {
	files, err := Getfiles(directory, opts)
	if err != nil {
		return err
	}

	fmt.Printf("%s:\n", directory)
	PrintFiles(files, opts)

	if opts.Recursive {
		for _, file := range files {
			if file.IsDir {
				ListDirectory(directory+"/"+file.Name, opts)
			}
		}
	}
	return nil
}


