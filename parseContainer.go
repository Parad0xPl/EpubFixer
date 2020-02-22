package main

import (
	"archive/zip"
	"encoding/xml"
)

type RootFiles struct {
	XMLName   xml.Name   `xml:"rootfiles"`
	RootFiles []RootFile `xml:"rootfile"`
}

type RootFile struct {
	XMLName   xml.Name `xml:"rootfile"`
	FullPath  string   `xml:"full-path,attr"`
	MediaType string   `xml:"media-type,attr"`
}

type ContainerXml struct {
	XMLName   xml.Name  `xml:"container"`
	Version   string    `xml:"version,attr"`
	Xmlns     string    `xml:"xmlns,attr"`
	RootFiles RootFiles `xml:"rootfiles"`
}

//ParseContainerFile parse container xml
func ParseContainerFile(f *zip.File) (*ContainerXml, error) {
	var err error
	reader, err := f.Open()
	if err != nil {
		return nil, err
	}
	xmlParser := xml.NewDecoder(reader)
	_, _ = xmlParser.Token()
	CTXml := new(ContainerXml)
	err = xmlParser.Decode(CTXml)
	if err != nil {
		return nil, err
	}
	return CTXml, err
}
