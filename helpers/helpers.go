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
			} else {
				fmt.Println("ls: option '--h' is ambiguous;")
				// break
				return Options{}, nil, errors.New("invalid option")
			}
			// should handle the error
			// ls: option '--h' is ambiguous; possibilities: '--human-readable' '--hide-control-chars' '--hide' '--hyperlink' '--help'
			// Try 'ls --help' for more information.
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
					return Options{}, nil, errors.New("invalid option")
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
		return nil, err
	}
	var files []FileInfo
	// var totalSize int64 = 0

	for _, file := range dir_entries {
		if !options.All && strings.HasPrefix(file.Name(), ".") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, err
		}

		status, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return nil, errors.New("failed to get file stats")
		}
		// Extract st_blocks instead of just size
		blocks := int64(status.Blocks) // st_blocks represents 512-byte blocks
		// lets extract the hardlinkd count
		links := status.Nlink

		user_info, err := user.LookupId(strconv.Itoa(int(status.Uid)))
		owner := "unknown"
		if err == nil {
			owner = user_info.Username
		}

		grp_info, err := user.LookupGroupId(strconv.Itoa(int(status.Gid)))
		group := "unknown"
		if err == nil {
			group = grp_info.Name
		}

		// totalSize += info.Size()
		files = append(files, FileInfo{
			Name:      info.Name(),
			Mode:      info.Mode(),
			ModTime:   info.ModTime(),
			Size:      info.Size(),
			Blocks:    blocks,
			IsDir:     info.IsDir(),
			Owner:     owner,
			Group:     group,
			HardLinks: int64(links),
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
			fmt.Printf("%s %d %s %s %4d %s %s\n", file.Mode, file.HardLinks, file.Owner, file.Group, file.Size, file.ModTime.Format("Jan 02 15:04"), file.Name)
		} else {
			fmt.Printf("%s  ", file.Name)
		}
		// fmt.Println()
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
	} else {
		// fmt.Printf("%s:\n", directory)
	}

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
