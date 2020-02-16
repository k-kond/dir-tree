package main

import (
	"fmt"
	"io"
	"os"
)

type levels map[int]bool

func printLine(writer io.Writer, file os.FileInfo, level levels) {
	sym := map[bool]string{
		true:  "└",
		false: "├",
	}

	symTab := map[bool]string{
		true:  "",
		false: "│",
	}
	//println(len(level))
	for i, j := range level {
		switch i {
		case len(level) - 1:
			fmt.Fprint(writer, sym[j]+"───")
		default:
			fmt.Fprint(writer, symTab[j]+"	")
		}
	}
	fmt.Fprintln(writer, file.Name())
}

func cOutTree(writer io.Writer, path string, printFiles bool, level levels) error {

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}

	files, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	cnt := len(level) - 1

	for i, file := range files {
		level[cnt] = i == len(files)-1
		if !file.IsDir() && !printFiles {
			continue
		} else {
			printLine(writer, file, level)
		}
		if file.IsDir() {
			level[cnt+1] = false
			cOutTree(writer, path+string(os.PathSeparator)+file.Name(), printFiles, level)
		}
	}
	return nil

}

func dirTree(writer io.Writer, path string, printFiles bool) error {

	lev := levels{
		0: false,
		//		1: false,
	}
	return cOutTree(writer, path, printFiles, lev)

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
