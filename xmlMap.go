// https://github.com/vikasgowda3007/golang

package main

import (
	"encoding/xml"
	"errors"
	"strings"
)

type xmlTestMap struct {
	Unparsed unparsedTagsMap `xml:",any"`
}

// unparsedTagMap contains the tag information
type unparsedTagMap struct {
	XMLName     xml.Name
	FullContent string `xml:",innerxml"`
}

// unparsedTagsMap store tags not handled by Unmarshal in a map, it should be labelled with `xml",any"`
type unparsedTagsMap map[string]string

var (
	finalMap  = make(map[string]string)
	separator string
)

// GetXmlMap accepts a xml string a the first parameter and a separator as the second parameter.
// Will return a map of xml data where recursive xml tag keys are separated by separator passed as parameter.
// If processing is unsuccessful xmlMap will be set to nil and err is returned.
func GetXmlMap(xmlStr string, sep string) (xmlMap map[string]string, err error) {
	var xmlStructMapMain xmlTestMap
	separator = sep
	err = xml.Unmarshal([]byte("<dummy>"+xmlStr+"</dummy>"), &xmlStructMapMain)
	if err != nil {
		return nil, err
	}
	if len(xmlStructMapMain.Unparsed) > 1 {
		return nil, errors.New("wrong XML input")
	}
	mainTag := ""
	for key := range xmlStructMapMain.Unparsed {
		mainTag = key
		break
	}
	if mainTag == "" {
		return nil, errors.New("wrong XML input : invalid tag name")
	}
	xmlStructMapMain.Unparsed = nil
	err = xml.Unmarshal([]byte(xmlStr), &xmlStructMapMain)
	if err != nil {
		return nil, err
	}
	err = recursive(xmlStructMapMain.Unparsed, mainTag)
	if err != nil {
		return nil, err
	}
	return finalMap, nil
}

func recursive(unparsed unparsedTagsMap, mainTag string) error {
	for key, value := range unparsed {
		if len(value) > 0 && key != "" {
			if strings.Contains(value, "<") {
				var xmlStructMapRec xmlTestMap
				err := xml.Unmarshal([]byte("<"+key+">"+value+"</"+key+">"), &xmlStructMapRec)
				if err != nil {
					return err
				}
				errRec := recursive(xmlStructMapRec.Unparsed, mainTag+separator+key)
				if errRec != nil {
					return errRec
				}
			} else {
				if mainTag != "" {
					finalMap[mainTag+separator+key] = value
				}
			}
		}
	}
	return nil
}

func (u *unparsedTagsMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if *u == nil {
		*u = unparsedTagsMap{}
	}
	e := unparsedTagMap{}
	err := d.DecodeElement(&e, &start)
	if err != nil {
		return err
	}

	(*u)[e.XMLName.Local] = e.FullContent
	return nil
}