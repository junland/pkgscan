package pkgscan

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// Listing describs a list of files.
type Listing struct {
	RealPath string
	RelPath  string
	Items    []FileInfo
}

// FileInfo describes a files attributes.
type FileInfo struct {
	Name      string
	Dir       bool
	Size      int64
	HumanSize string
	LastMod   string
}

type sortDirNameFirst Listing

func main() {
	list, err := DirList("/home/john")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(list)
}

func (l sortDirNameFirst) Len() int      { return len(l.Items) }
func (l sortDirNameFirst) Swap(i, j int) { l.Items[i], l.Items[j] = l.Items[j], l.Items[i] }

// Dont worry about cases.
func (l sortDirNameFirst) Less(i, j int) bool {
	return l.Items[i].Dir
}

// DirList lists the information about files inside a directory.
func DirList(file string) (Listing, error) {
	var list []FileInfo

	fmt.Printf("Looking up dir: %s", file)
	f, err := os.Open(file)
	if err != nil {
		return Listing{}, err
	}
	fi, err := f.Stat()
	if err != nil {
		return Listing{}, err
	}
	defer f.Close()

	if fi.IsDir() {
		files, err := ioutil.ReadDir(file)
		if err != nil {
			return Listing{}, err
		}

		// Start going thru each file and do stuff.
		for _, f := range files {

			// file name
			name := f.Name()
			if f.IsDir() {
				name += "/"
			}

			// skip hidden files.
			if strings.HasPrefix(name, ".") {
				continue
			}

			// file type
			dir := f.IsDir()

			// file size
			size := f.Size()

			// human file size
			hsize := ByteCountBinary(size)

			// file last mod time
			mod := f.ModTime().Format("2006-01-02 15:04")

			list = append(list, FileInfo{Name: name, Dir: dir, Size: size, HumanSize: hsize, LastMod: mod})
		}
		fmt.Println("Before sort: ", list)
		sort.Sort(sortDirNameFirst(Listing{RealPath: file, Items: list}))
		fmt.Println("After sort: ", list)
		return Listing{RealPath: file, Items: list}, nil
	}

	return Listing{}, fmt.Errorf("%s is not a directory", fi)
}