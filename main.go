package main

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type Metadata struct {
	Title   string   `xml:"title" json:"title"`
	Creator []string `xml:"creator" json:"creator"`
}

type Package struct {
	Metadata Metadata `xml:"metadata"`
}

type Container struct {
	Rootfiles struct {
		Rootfile struct {
			FullPath string `xml:"full-path,attr"`
		} `xml:"rootfile"`
	} `xml:"rootfiles"`
}

type Ebook struct {
	FileName string `json:"fileName"`
	Metadata
}

type Collection struct {
	Ebooks []Ebook `json:"ebooks"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <epub_file>...", os.Args[0])
	}

	collection := Collection{Ebooks: make([]Ebook, 0)}
	for _, fileName := range os.Args[1:] {
		ebook, err := read(fileName)
		if err != nil {
			log.Printf("reading %s: %s", fileName, err)
			continue
		}
		collection.Ebooks = append(collection.Ebooks, ebook)
	}

	b, err := json.Marshal(collection)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func read(fileName string) (Ebook, error) {
	zipFile, err := zip.OpenReader(fileName)
	if err != nil {
		return Ebook{}, err
	}
	defer zipFile.Close()

	var container Container
	if err := unmarshal(zipFile, "META-INF/container.xml", &container); err != nil {
		return Ebook{}, err
	}

	var pkg Package
	if err := unmarshal(zipFile, container.Rootfiles.Rootfile.FullPath, &pkg); err != nil {
		return Ebook{}, err
	}

	return Ebook{FileName: fileName, Metadata: pkg.Metadata}, nil
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
