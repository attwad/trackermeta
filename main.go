package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/attwad/trackermeta/it"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	paths, err := filepath.Glob("/Users/tmu/Downloads/Bogdan/*.it")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range paths {
		fmt.Println(it.ReadITFile(f))
	}
	fmt.Println("Read", len(paths), "files")
}
