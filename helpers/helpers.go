// helpers/helpers.go
package helpers

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var Help_msg = `Usage: my-ls [OPTION]... [FILE]...
List information about the FILEs (the current directory by default).

Mandatory arguments to long options are mandatory for short options too.
  -a, --all                  do not ignore entries starting with .
  -l                         use a long listing format
  -r, --reverse              reverse order while sorting
  -R, --recursive            list subdirectories recursively
  -t                         sort by time, newest first; see --time`

type FileInfo struct {
	Name      string
	Mode      fs.FileMode
	Size      int64
	Blocks    int64 // Store block count
	ModTime   time.Time
	IsDir     bool
	Owner     string
	Group     string
	HardLinks int64
}

type Options struct {
	Long      bool
	Recursive bool
	All       bool
	Reverse   bool
	TimeSort  bool
}

func ParseFlags(args []string) (Options, []string, error) {
	opts := Options{}
	var dirs []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			if arg == "--rev" || arg == "--reverse" {
				opts.Reverse = true
			} else if arg == "--rec" || arg == "--recursive" {
				opts.Recursive = true
			} else if arg == "--all " {
				opts.All = true
			} else if arg == "--help" {
				return Options{}, nil, errors.New(Help_msg)
			} else {
				return Options{}, nil, fmt.Errorf("my-ls: option '%s' is ambiguous.\nTry 'my-ls --help' for more information", arg)
			}
		} else if strings.HasPrefix(arg, "-") {
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
				default:
					return Options{}, nil, fmt.Errorf("my-ls: invalid option -- '%s'\nTry 'my-ls --help' for more information", arg)
				}
			}
		} else {
			dirs = append(dirs, arg)
		}
	}

	if len(dirs) == 0 {
		dirs = append(dirs, ".")
	}

	return opts, dirs, nil
}

func GetFiles(directory string, options Options) ([]FileInfo, error) {
	dir_entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("my-ls: cannot access '%s': No such file or directory", directory)
	}

	var files []FileInfo
	// Include "." and ".." directories
	dotFiles := []string{".", ".."}
	for _, name := range dotFiles {
		info, err := os.Stat(directory + "/" + name)
		if err == nil { // Only add if accessible
			files = append(files, getFileInfo(info))
		}
	}

	for _, file := range dir_entries {
		if !options.All && strings.HasPrefix(file.Name(), ".") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, getFileInfo(info))
	}

	SortFiles(&files, options)
	return files, nil
}

func SortFiles(files *[]FileInfo, opts Options) {
	if opts.TimeSort {
		Sort_by_Time(files)
	}
	if opts.Reverse {
		Sort_by_reverse(files)
	}
}

func PrintFiles(files []FileInfo, opts Options) {
	for _, file := range files {
		if !opts.All && (file.Name == "." || file.Name == "..") {
			continue
		}
		if opts.Long {
			fmt.Printf("%s %d %s %s %4d %s %s\n", file.Mode, file.HardLinks, file.Owner, file.Group, file.Size, file.ModTime.Format("Jan 02 15:04"), file.Name)
		} else {
			fmt.Printf("%s  ", file.Name)
		}

	}

	fmt.Println()
}

func ListDirectory(directory string, opts Options) error {
	files, err := GetFiles(directory, opts)
	if err != nil {
		return err
	}
	// Calculate total block size
	var totalBlocks int64
	for _, file := range files {
		totalBlocks += file.Blocks
	}
	if opts.Long {
		fmt.Printf("%s:\ntotal %d\n", directory, totalBlocks/2) // Convert blocks to 1024-byte units
	}
	if !opts.Long && opts.Recursive {
		fmt.Printf("%s:\n", directory)
	}

	PrintFiles(files, opts)

	if opts.Recursive {
		for _, file := range files {
			if file.IsDir && file.Name != "." && file.Name != ".." {
				ListDirectory(directory+"/"+file.Name, opts)
			}
		}
	}
	return nil
}

// Helper function to extract file info
func getFileInfo(info fs.FileInfo) FileInfo {
	status, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic("failed to get file stats") // Should never happen
	}

	blocks := int64(status.Blocks) // 512-byte blocks
	links := status.Nlink

	owner, group := "unknown", "unknown"
	if userInfo, err := user.LookupId(strconv.Itoa(int(status.Uid))); err == nil {
		owner = userInfo.Username
	}
	if grpInfo, err := user.LookupGroupId(strconv.Itoa(int(status.Gid))); err == nil {
		group = grpInfo.Name
	}

	return FileInfo{
		Name:      info.Name(),
		Mode:      info.Mode(),
		ModTime:   info.ModTime(),
		Size:      info.Size(),
		Blocks:    blocks,
		IsDir:     info.IsDir(),
		Owner:     owner,
		Group:     group,
		HardLinks: int64(links),
	}
}

func Sort_by_Time(files *[]FileInfo) {
	for i := 0; i < len(*files); i++ {
		for j := i + 1; j < len(*files); j++ {
			time1 := (*files)[i].ModTime
			time2 := (*files)[j].ModTime
			if time2.After(time1) {
				(*files)[i], (*files)[j] = (*files)[j], (*files)[i]
			}
		}
	}
}

func Sort_by_reverse(files *[]FileInfo) {
	for i := 0; i < len(*files); i++ {
		for j := i + 1; j < len(*files); j++ {
			name1 := (*files)[i].Name
			name2 := (*files)[j].Name
			if name1 < name2 {
				(*files)[i], (*files)[j] = (*files)[j], (*files)[i]
			}
		}
	}
}
