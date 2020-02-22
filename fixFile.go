package main

import (
	"archive/zip"
	"fmt"
	"os"
)

//ReadAllFileFromZip Read all data from and return it as string
func ReadAllFileFromZip(zipFile *zip.File) (string, error) {
	var err error
	reader, err := zipFile.Open()
	if err != nil {
		return "", err
	}
	// Possible error when platform is not 64 bit
	buffer := make([]byte, zipFile.UncompressedSize64)
	_, err = reader.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer), err
}

//ProcessFile Perform fix on file
func ProcessFile(inputFilename string) error {
	var err error
	outputFilename := inputFilename + ".fixed.epub"
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer inputFile.Close()
	inputStats, err := inputFile.Stat()
	if err != nil {
		return err
	}
	zipReader, err := zip.NewReader(inputFile, inputStats.Size())
	if err != nil {
		return err
	}

	book, err := ParseEPUB(zipReader)
	if err != nil {
		return err
	}

	PatchBook(book)

	err = SaveEPUB(book, outputFilename)
	if err != nil {
		return err
	}

	return err
}

//PatchBook modify book structure with fixes
func PatchBook(book *EPUB) {
	for _, sr := range book.RootFiles {
		ui := sr.UniqueIdentifier
		var newIdentifiers []DCIdentifierUnmarshal
		for _, el := range sr.MetaData.(*MetaDataUnmarshal).Identifier {
			if el.ID == ui {
				newIdentifiers = []DCIdentifierUnmarshal{el}
				break
			}
		}
		sr.MetaData.(*MetaDataUnmarshal).Identifier = newIdentifiers

		newItems := make([]Item, 0)
		for _, el := range sr.Manifest.Items {
			if el.ID == "default-info" {
				fmt.Println("Skipping")
			} else {
				newItems = append(newItems, el)
			}
		}
		sr.Manifest.Items = newItems

		newRefs := make([]ItemRef, 0)
		for _, el := range sr.Spine.ItemRefs {
			if el.IDRef == "default-info" {
				fmt.Println("Skipping")
			} else {
				newRefs = append(newRefs, el)
			}
		}
		sr.Spine.ItemRefs = newRefs
	}

}
