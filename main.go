package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type levels map[int]bool

var sym = map[bool]string{
	true:  "└",
	false: "├",
}
var symTab = map[bool]string{
	true:  "",
	false: "│",
}

func printLine(writer io.Writer, file *os.FileInfo, levelMap *levels, level int, pFiles bool) {
	out := []string{}
	for i := 0; i <= level; i++ {
		statMap := (*levelMap)[i]
		switch i {
		case level:
			out = append(out, sym[statMap], "───")
		default:
			out = append(out, symTab[statMap], "\t")
		}

	}
	out = append(out, (*file).Name())
	if pFiles && !(*file).IsDir() {
		fSize := (*file).Size()
		if fSize == 0 {
			out = append(out, " (empty)")
		} else {
			out = append(out, " (", strconv.FormatInt(fSize, 10), "b)")
		}
	}
	fmt.Fprintln(writer, strings.Join(out, ""))
}

func onlyToWant(files *[]os.FileInfo, printFiles bool) {

	if !printFiles {
		for i := len(*files) - 1; i >= 0; i-- {
			if !(*files)[i].IsDir() {
				*files = append((*files)[:i], (*files)[i+1:]...)
			}
		}
	}
}

func cOutTree(writer io.Writer, path string, printFiles bool, levelMap *levels, itLevel int) error {

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}

	files, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	onlyToWant(&files, printFiles)
	for i, file := range files {
		fileName := file.Name()
		(*levelMap)[itLevel] = i == (len(files) - 1) // end of list?
		printLine(writer, &file, levelMap, itLevel, printFiles)

		if file.IsDir() {
			cOutTree(writer, path+string(os.PathSeparator)+fileName, printFiles, levelMap, itLevel+1)
		}
	}
	return nil
}

func dirTree(writer io.Writer, path string, printFiles bool) error {

	lev := levels{
		0: false,
	}
	//fmt.Fprintln(writer, path)
	return cOutTree(writer, path, printFiles, &lev, 0)

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
