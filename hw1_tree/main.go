package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func size(s int64) string {
	if s == 0 {
		return "empty"
	} else {
		return strconv.FormatInt(s, 10) + "b"
	}
}

func readDir(dirname string, printFiles bool) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	if printFiles {
		return list, nil
	} else {
		var dirs []os.FileInfo
		for _, v := range list {
			if v.IsDir() {
				dirs = append(dirs, v)
			}
		}
		return dirs, nil
	}
}

func walk(out io.Writer, root, prev string, printFiles bool) error {
	dir, err := readDir(root, printFiles)
	if err != nil {
		return err
	}

	for i, f := range dir {
		next := "│	"
		if i == len(dir)-1 {
			fmt.Fprint(out, prev+"└───")
			next = "	"
		} else {
			fmt.Fprint(out, prev+"├───")
		}

		if f.IsDir() {
			fmt.Fprintln(out, f.Name())

			err := walk(out, filepath.Join(root, f.Name()), prev+next, printFiles)
			if err != nil {
				return err
			}

		} else if printFiles {
			fmt.Fprintf(out, "%s (%s)\n", f.Name(), size(f.Size()))
		}
	}

	return nil
}

func dirTree(out io.Writer, root string, printFiles bool) error {
	r, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("%s directory does not exist: %v", root, err)
	}

	if !r.IsDir() {
		return fmt.Errorf("%s is not directory", root)
	}

	err = walk(out, root, "", printFiles)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
