package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/attwad/trackermeta/data"
	"github.com/attwad/trackermeta/html"
	"github.com/attwad/trackermeta/it"
	"github.com/attwad/trackermeta/xm"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	paths, err := filepath.Glob("/Users/tmu/Downloads/Bogdan/*")
	if err != nil {
		log.Fatal(err)
	}
	m := make([]data.TrackerFile, 0)
	for _, f := range paths {
		if strings.HasSuffix(f, ".it") {
			tm, err := it.ReadITFile(f)
			if err != nil {
				log.Fatal(err)
			}
			m = append(m, *tm)
		} else if strings.HasSuffix(f, ".xm") {
			tm, err := xm.ReadXMFile(f)
			if err != nil {
				log.Fatal(err)
			}
			m = append(m, *tm)
		}
	}
	fmt.Println("Read", len(paths), "files")
	f, err := os.Create("out.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := html.RenderHTML(m, f); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote out.html")
}
