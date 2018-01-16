package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-fs"
	"github.com/mitchellh/go-fs/fat"
)

func main() {
	args := os.Args[1:]
	nargs := len(args)
	if nargs < 2 {
		fmt.Println("Usage: fatcopy file [file]... image")
		return
	}
	imgFile := args[nargs-1]
	srcFiles := args[:nargs-1]
	f, err := os.OpenFile(imgFile, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	device, err := fs.NewFileDisk(f)
	if err != nil {
		panic(err)
	}

	filesys, err := fat.New(device)
	if err != nil {
		panic(err)
	}

	rootDir, err := filesys.RootDir()
	if err != nil {
		panic(err)
	}

	for _, fname := range srcFiles {
		subEntry, err := rootDir.AddFile(filepath.Base(fname))
		if err != nil {
			panic(err)
		}

		file, err := subEntry.File()
		if err != nil {
			panic(err)
		}

		src, err := os.Open(fname)
		if err != nil {
			log.Printf("Unable to open file %s: %v", fname, err)
		}

		_, err = io.Copy(file, src)
		if err != nil {
			panic(err)
		}
	}
}
