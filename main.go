package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type Package struct {
	Metadata struct {
		Title   string   `xml:"title"`
		Creator []string `xml:"creator"`
	} `xml:"metadata"`
}

type Container struct {
	Rootfiles struct {
		Rootfile struct {
			FullPath string `xml:"full-path,attr"`
		} `xml:"rootfile"`
	} `xml:"rootfiles"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <epub_file>", os.Args[0])
	}

	zipFile, err := zip.OpenReader(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	var container Container

	for _, f := range zipFile.File {
		if f.Name == "META-INF/container.xml" {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer rc.Close()

			if err := xml.NewDecoder(rc).Decode(&container); err != nil {
				log.Fatal(err)
			}

			break
		}
	}

	var pkg Package

	for _, f := range zipFile.File {
		if f.Name == container.Rootfiles.Rootfile.FullPath {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer rc.Close()

			if err := xml.NewDecoder(rc).Decode(&pkg); err != nil {
				log.Fatal(err)
			}

			break
		}
	}

	fmt.Printf("Title: %s\n", pkg.Metadata.Title)
	fmt.Print("Author(s): ")
	for i, author := range pkg.Metadata.Creator {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s", author)
	}
	fmt.Println()
}
