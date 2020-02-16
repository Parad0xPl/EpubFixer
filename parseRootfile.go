package main

import (
	"archive/zip"
	"encoding/xml"
)

type DCTitleUnmarshal struct {
	XMLName xml.Name `xml:"http://purl.org/dc/elements/1.1/ title"`
	Title   string   `xml:",chardata"`
}

type DCTitleMarshal struct {
	XMLName xml.Name `xml:"dc:title"`
	Title   string   `xml:",chardata"`
}

type DCCreatorUnmarshal struct {
	XMLName xml.Name `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Creator string   `xml:",chardata"`
	Role    string   `xml:"http://www.idpf.org/2007/opf role,attr"`
	FileAs  string   `xml:"http://www.idpf.org/2007/opf file-as,attr"`
}

type DCCreatorMarshal struct {
	XMLName xml.Name `xml:"dc:creator"`
	Creator string   `xml:",chardata"`
	Role    string   `xml:"opf:role,attr"`
	FileAs  string   `xml:"opf:file-as,attr"`
}

type DCPublisherUnmarshal struct {
	XMLName   xml.Name `xml:"http://purl.org/dc/elements/1.1/ publisher"`
	Publisher string   `xml:",chardata"`
}

type DCPublisherMarshal struct {
	XMLName   xml.Name `xml:"dc:publisher"`
	Publisher string   `xml:",chardata"`
}

type DCIdentifierUnmarshal struct {
	XMLName    xml.Name `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	Identifier string   `xml:",chardata"`
	ID         string   `xml:"id,attr"`
}

type DCIdentifierMarshal struct {
	XMLName    xml.Name `xml:"dc:identifier"`
	Identifier string   `xml:",chardata"`
	ID         string   `xml:"id,attr"`
}

type DCLanguageUnmarshal struct {
	XMLName  xml.Name `xml:"http://purl.org/dc/elements/1.1/ language"`
	Language string   `xml:",chardata"`
}

type DCLanguageMarshal struct {
	XMLName  xml.Name `xml:"dc:language"`
	Language string   `xml:",chardata"`
}

type Meta struct {
	XMLName xml.Name `xml:"meta"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:"content,attr"`
}

type MetaDataUnmarshal struct {
	XMLName    xml.Name                `xml:"metadata"`
	XmlnsDc    string                  `xml:"xmlns:dc,attr"`
	XmlnsOpf   string                  `xml:"xmlns:opf,attr"`
	Title      DCTitleUnmarshal        `xml:"http://purl.org/dc/elements/1.1/ title"`
	Creator    DCCreatorUnmarshal      `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Publisher  DCPublisherUnmarshal    `xml:"http://purl.org/dc/elements/1.1/ publisher"`
	Identifier []DCIdentifierUnmarshal `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	Language   DCLanguageUnmarshal     `xml:"http://purl.org/dc/elements/1.1/ language"`
	Metas      []Meta                  `xml:"meta"`
}

type MetaDataMarshal struct {
	XMLName    xml.Name              `xml:"metadata"`
	XmlnsDc    string                `xml:"xmlns:dc,attr"`
	XmlnsOpf   string                `xml:"xmlns:opf,attr"`
	Title      DCTitleMarshal        `xml:"dc:title"`
	Creator    DCCreatorMarshal      `xml:"dc:creator"`
	Publisher  DCPublisherMarshal    `xml:"dc:publisher"`
	Identifier []DCIdentifierMarshal `xml:"dc:identifier"`
	Language   DCLanguageMarshal     `xml:"dc:language"`
	Metas      []Meta                `xml:"meta"`
}

func TransformFromUnmarshal(old *MetaDataUnmarshal) *MetaDataMarshal {
	identifiersLen := len(old.Identifier)
	transofrmed := &MetaDataMarshal{
		XMLName:  old.XMLName,
		XmlnsDc:  old.XmlnsDc,
		XmlnsOpf: old.XmlnsOpf,
		Title: DCTitleMarshal{
			XMLName: old.Title.XMLName,
			Title:   old.Title.Title,
		},
		Creator: DCCreatorMarshal{
			XMLName: old.Creator.XMLName,
			Creator: old.Creator.Creator,
			FileAs:  old.Creator.FileAs,
			Role:    old.Creator.Role,
		},
		Publisher: DCPublisherMarshal{
			XMLName:   old.Publisher.XMLName,
			Publisher: old.Publisher.Publisher,
		},
		Identifier: make([]DCIdentifierMarshal, identifiersLen),
		Language: DCLanguageMarshal{
			XMLName:  old.Language.XMLName,
			Language: old.Language.Language,
		},
		Metas: old.Metas,
	}

	for i, el := range old.Identifier {
		transofrmed.Identifier[i] = DCIdentifierMarshal{
			XMLName:    el.XMLName,
			Identifier: el.Identifier,
			ID:         el.ID,
		}
	}

	return transofrmed
}

type Item struct {
	XMLName   xml.Name `xml:"item"`
	ID        string   `xml:"id,attr"`
	Href      string   `xml:"href,attr"`
	MediaType string   `xml:"media-type,attr"`
}

type Manifest struct {
	XMLName xml.Name `xml:"manifest"`
	Items   []Item   `xml:"item"`
}

type ItemRef struct {
	XMLName xml.Name `xml:"itemref"`
	IDRef   string   `xml:"idref,attr"`
}

type Spine struct {
	XMLName  xml.Name  `xml:"spine"`
	TOC      string    `xml:"toc,attr"`
	ItemRefs []ItemRef `xml:"itemref"`
}

type Reference struct {
	XMLName xml.Name `xml:"reference"`
	Href    string   `xml:"href,attr"`
	Title   string   `xml:"title,attr"`
	Type    string   `xml:"type,attr"`
}

type Guide struct {
	XMLName    xml.Name    `xml:"guide"`
	References []Reference `xml:"reference"`
}

type RootFileXml struct {
	XMLName          xml.Name    `xml:"package"`
	Xmlns            string      `xml:"xmlns,attr"`
	Version          string      `xml:"version,attr"`
	UniqueIdentifier string      `xml:"unique-identifier,attr"`
	MetaData         interface{} `xml:"metadata"`
	Manifest         Manifest    `xml:"manifest"`
	Spine            Spine       `xml:"spine"`
	Guide            Guide       `xml:"guide"`
}

func newUnmarshal() *RootFileXml {
	obj := new(RootFileXml)
	obj.MetaData = new(MetaDataUnmarshal)
	return obj
}

func newMarshal() *RootFileXml {
	obj := new(RootFileXml)
	obj.MetaData = new(MetaDataMarshal)
	return obj
}

func ParseRootFile(f *zip.File) (*RootFileXml, error) {
	reader, err := f.Open()
	if err != nil {
		return nil, err
	}
	xmlParser := xml.NewDecoder(reader)
	_, _ = xmlParser.Token()
	RFXml := newUnmarshal()
	err = xmlParser.Decode(RFXml)
	if err != nil {
		return nil, err
	}
	return RFXml, nil
}
