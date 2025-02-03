package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func formatDirEntry(entry os.DirEntry, printFiles bool, hasNoBroMas []bool) string {
	var indent string
	prefix := "├───"
	if hasNoBroMas[len(hasNoBroMas)-1] {
		prefix = "└───"
	}
	for _, elem := range hasNoBroMas[:len(hasNoBroMas)-1] {
		if !elem {
			indent += "│" + "\t"
		} else {
			indent += "\t"
		}
	}
	resStr := fmt.Sprintf("%s%s%s", indent, prefix, entry.Name())
	if printFiles {
		if !entry.IsDir() {
			fileInfo, _ := entry.Info()
			size := fileInfo.Size()
			if size == 0 {
				resStr = fmt.Sprintf("%s (empty)", resStr)
			} else {
				resStr = fmt.Sprintf("%s (%db)", resStr, size)
			}
		}
	}
	resStr += "\n"
	return resStr
}

func recursiveDirTree(output io.Writer, path string, printFiles bool, level int, hasNoBroMas []bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	hasNoBroMas = hasNoBroMas[:level]

	sortByIncreasing := func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	}
	sort.Slice(entries, sortByIncreasing)
	filtredEntries := make([]os.DirEntry, 0)
	if !printFiles {
		for _, entry := range entries {
			if entry.IsDir() {
				filtredEntries = append(filtredEntries, entry)
			}
		}
	} else {
		filtredEntries = entries
	}

	for i, entry := range filtredEntries {
		isLast := i == len(filtredEntries)-1
		if len(hasNoBroMas)-1 < level {
			hasNoBroMas = append(hasNoBroMas, isLast)
		}
		if level == len(hasNoBroMas)-1 && isLast {
			hasNoBroMas[level] = true
		}
		if !entry.IsDir() && printFiles {
			_, err = fmt.Fprintf(output, formatDirEntry(entry, printFiles, hasNoBroMas))
			if err != nil {
				return err
			}
		}
		if entry.IsDir() {
			_, err = fmt.Fprintf(output, formatDirEntry(entry, printFiles, hasNoBroMas))
			if err != nil {
				return err
			}
			err = recursiveDirTree(output, filepath.Join(path, entry.Name()), printFiles, level+1, hasNoBroMas)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	return recursiveDirTree(output, path, printFiles, 0, make([]bool, 0))
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
