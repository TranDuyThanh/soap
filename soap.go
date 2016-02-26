package soap

import (
	// "fmt"
	"strings"
)

var code = map[string]string{
	"&":  "&amp;",
	"<":  "&lt;",
	">":  "&gt;",
	"'":  "&apos;",
	"\"": "&quot;",
}

func EncodeXML(xmlString string) string {
	xmlString = strings.Replace(xmlString, "\n", "", -1)
	xmlString = strings.Replace(xmlString, " ", "", -1)
	for key, value := range code {
		xmlString = strings.Replace(xmlString, key, value, -1)
	}
	return xmlString
}

func DecodeXML(xmlString string) string {
	for key, value := range code {
		xmlString = strings.Replace(xmlString, value, key, -1)
	}
	return xmlString
}
