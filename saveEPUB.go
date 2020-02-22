package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

//WriteMimeType create "mimetype" file in archive
func WriteMimeType(zipWriter *zip.Writer) error {
	file, err := zipWriter.Create("mimetype")
	if err != nil {
		return err
	}
	_, err = file.Write([]byte("application/epub+zip"))
	if err != nil {
		return err
	}
	return nil
}

//WriteContainer serialize and write container file in archive
func WriteContainer(book *EPUB, zipWriter *zip.Writer) error {
	file, err := zipWriter.Create("META-INF/container.xml")
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	err = encoder.Encode(book.Container)
	if err != nil {
		return err
	}

	return nil
}

//WriteRootFiles serialize and write root files of epub
func WriteRootFiles(book *EPUB, zipWriter *zip.Writer) error {
	for rfPath, rf := range book.RootFiles {
		rf.MetaData.(*MetaDataUnmarshal).XmlnsDc = "http://purl.org/dc/elements/1.1/"
		rf.MetaData.(*MetaDataUnmarshal).XmlnsOpf = "http://www.idpf.org/2007/opf"

		fmt.Println(rf.MetaData.(*MetaDataUnmarshal))

		rf.MetaData = TransformFromUnmarshal(rf.MetaData.(*MetaDataUnmarshal))

		fmt.Println(rf.MetaData.(*MetaDataMarshal))

		file, err := zipWriter.Create(rfPath)
		if err != nil {
			return err
		}
		_, err = file.Write([]byte(xml.Header))
		if err != nil {
			return err
		}
		enc := xml.NewEncoder(file)
		enc.Indent("", "  ")
		err = enc.Encode(rf)
		if err != nil {
			return err
		}

	}
	return nil
}

//WriteRefFiles save all referenced files from each root file to archive
func WriteRefFiles(book *EPUB, w *zip.Writer) error {
	for rfPath, rf := range book.RootFiles {
		dir := path.Dir(rfPath)
		for _, file := range rf.Manifest.Items {
			targetPath := path.Join(dir, file.Href)
			bookFile, ok := book.Files[targetPath]
			if !ok {
				fmt.Printf("Can't find refereed file '%s' in structure cache. RootFile path: '%s'\n", targetPath, rfPath)
				continue
			}
			file, err := w.Create(targetPath)
			if err != nil {
				return err
			}
			bookFileReader, err := bookFile.Open()
			if err != nil {
				return err
			}
			buffer, err := ioutil.ReadAll(bookFileReader)
			if err != nil {
				return err
			}
			_, err = file.Write(buffer)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//SaveEPUB serialize structures and write it to file
func SaveEPUB(book *EPUB, filename string) error {
	var err error

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()
	err = WriteMimeType(w)
	if err != nil {
		return err
	}
	err = WriteContainer(book, w)
	if err != nil {
		return err
	}
	err = WriteRootFiles(book, w)
	if err != nil {
		return err
	}
	err = WriteRefFiles(book, w)
	if err != nil {
		return err
	}

	return err
}
