package main

import (
	"archive/zip"
	"fmt"
)

//EPUB holds structured data needed for recreating file
type EPUB struct {
	MimeType  string
	Files     map[string]*zip.File
	Container *ContainerXml
	RootFiles map[string]*RootFileXml
}

//ParseEPUB reads zip/epub file and bind content to EPUB struct
func ParseEPUB(zipReader *zip.Reader) (*EPUB, error) {
	book := new(EPUB)
	book.Files = make(map[string]*zip.File)
	book.RootFiles = make(map[string]*RootFileXml)
	for _, f := range zipReader.File {
		book.Files[f.Name] = f
		if f.Name == "mimetype" {
			fmt.Printf("Found mimeType\n")
			mimeType, err := ReadAllFileFromZip(f)
			if err != nil {
				return nil, err
			}
			book.MimeType = mimeType
		} else if f.Name == "META-INF/container.xml" {
			ct, err := ParseContainerFile(f)
			if err != nil {
				return nil, err
			}
			book.Container = ct
		}
	}

	if book.Container != nil {
		for _, rootFile := range book.Container.RootFiles.RootFiles {
			fmt.Println("Parsing", rootFile.FullPath)
			file, ok := book.Files[rootFile.FullPath]
			if ok {
				rf, err := ParseRootFile(file)
				if err != nil {
					return nil, err
				}
				book.RootFiles[rootFile.FullPath] = rf
			}
		}
	}
	return book, nil
}
