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
	if err := unmarshal(zipFile, "META-INF/container.xml", &container); err != nil {
		log.Fatal(err)
	}

	var pkg Package
	if err := unmarshal(zipFile, container.Rootfiles.Rootfile.FullPath, &pkg); err != nil {
		log.Fatal(err)
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

func unmarshal(zipFile *zip.ReadCloser, fullPath string, v any) error {
	for _, f := range zipFile.File {
		if f.Name == fullPath {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("error reading %s: %w", fullPath, err)
			}
			defer rc.Close()

			if err := xml.NewDecoder(rc).Decode(v); err != nil {
				return fmt.Errorf("error unmarshaling %s: %w", fullPath, err)
			}

			return nil
		}
	}
	return fmt.Errorf("error reading %s: file not found", fullPath)
}
