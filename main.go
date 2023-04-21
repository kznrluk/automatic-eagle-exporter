package main

import (
	"github.com/kznrluk/sdweb-eaglepack/export"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <GLOB pattern>", os.Args[0])
	}

	pattern := os.Args[1]

	matchingFiles, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("Error searching for files: %v", err)
	}

	var loaded []*export.ExportImage
	for _, file := range matchingFiles {
		image, err := export.CreateExportImage(file)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		log.Printf("Load %s", file)
		loaded = append(loaded, image)
	}

	log.Printf("Create eaglepack...")
	err = export.CreateZip(loaded)
	if err != nil {
		panic(err)
	}
	log.Printf("Done!")
}
