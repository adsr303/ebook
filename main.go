package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <epub_file>", os.Args[0])
	}

	zipFile, err := zip.OpenReader(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	for _, f := range zipFile.File {
		fmt.Println(f.Name)
	}
}
