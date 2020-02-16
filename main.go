package main

import (
	"fmt"
	"io"
	"os"
)

type levels map[int]bool

func printLine(writer io.Writer, fileName string, levelMap levels, level int) {
	sym := map[bool]string{
		true:  "└",
		false: "├",
	}

	symTab := map[bool]string{
		true:  "",
		false: "│",
	}
	//println(len(levelMap))
	for i := 0; i < level; i++ {
		switch i {
		case level - 1:
			fmt.Fprint(writer, sym[levelMap[i]]+"───")
		default:
			//fmt.Fprint(writer, symTab[levelMap[i]]+"	")
			if levelMap[i] {
				fmt.Fprint(writer, symTab[levelMap[i]]+"───")
			} else {
				fmt.Fprint(writer, symTab[levelMap[i]]+"	")
			}

		}

	}

	fmt.Fprintln(writer, fileName)
}

func onlyToWant(files *[]os.FileInfo, printFiles bool) {

	if !printFiles {
		for i := len(*files) - 1; i >= 0; i-- {
			// name := (*files)[i].Name()
			// isdir := (*files)[i].IsDir()
			// fmt.Println(name)
			if !(*files)[i].IsDir() {
				*files = append((*files)[:i], (*files)[i+1:]...)
			}
		}
	}
}

func cOutTree(writer io.Writer, path string, printFiles bool, levelMap levels, itLevel int) error {

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
	for _, file := range files {
		fileName := file.Name()
		//levelMap[itLevel] = i == (len(files) - 1)
		printLine(writer, fileName, levelMap, itLevel)

		if file.IsDir() {
			levelMap[itLevel] = true
			levelMap[itLevel+1] = false
			cOutTree(writer, path+string(os.PathSeparator)+fileName, printFiles, levelMap, itLevel+1)
		}
	}
	return nil

}

func dirTree(writer io.Writer, path string, printFiles bool) error {

	lev := levels{
		0: false,
		//		1: false,
	}
	printLine(writer, path, lev, 0)
	return cOutTree(writer, path, printFiles, lev, 1)

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
