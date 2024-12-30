package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type Metadata struct {
	XMLName  xml.Name `xml:"package"`
	Metadata struct {
		Title   string   `xml:"title"`
		Creator []string `xml:"creator"`
	} `xml:"metadata"`
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

	var metadata Metadata

	var container struct {
		Rootfiles struct {
			Rootfile struct {
				FullPath string `xml:"full-path,attr"`
			} `xml:"rootfile"`
		} `xml:"rootfiles"`
	}

	for _, f := range zipFile.File {
		if f.Name == "META-INF/container.xml" {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer rc.Close()

			// Parse container.xml to find the path to the OPF file
			if err := xml.NewDecoder(rc).Decode(&container); err != nil {
				log.Fatal(err)
			}

			break
		}
	}

	// Find and parse the OPF file
	for _, f := range zipFile.File {
		if f.Name == container.Rootfiles.Rootfile.FullPath {
			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer rc.Close()

			if err := xml.NewDecoder(rc).Decode(&metadata); err != nil {
				log.Fatal(err)
			}

			break
		}
	}

	fmt.Printf("Title: %s\n", metadata.Metadata.Title)
	fmt.Printf("Author(s): ")
	for i, author := range metadata.Metadata.Creator {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", author)
	}
	fmt.Println()
}
