package soap

import (
	// "fmt"
	"regexp"
	"strings"
)

var code = map[string]string{
	"&":  "&amp;",
	"<":  "&lt;",
	">":  "&gt;",
	"'":  "&apos;",
	"\"": "&quot;",
}

func EscapeXML(xmlString string) string {
	xmlString = strings.Replace(xmlString, "\n", "", -1)
	xmlString = strings.Replace(xmlString, " ", "", -1)
	for key, value := range code {
		xmlString = strings.Replace(xmlString, key, value, -1)
	}
	return xmlString
}

func UnescapeXML(xmlString string) string {
	for key, value := range code {
		xmlString = strings.Replace(xmlString, value, key, -1)
	}
	return xmlString
}

func GetResponse(xmlString string) string {
	r, _ := regexp.Compile("(&lt;resultInfo)(.+)(resultInfo&gt;)")
	newString := r.FindString(xmlString)
	return newString
}

func GetUnescapedResponse(xmlString string) string {
	xmlString = GetResponse(xmlString)
	newString := UnescapeXML(xmlString)
	return newString
}

func InitPayload(template, action, condition string) string {
	newString := strings.Replace(template, "*ACTION*", action, -1)
	newString = strings.Replace(newString, "*CONDITION*", condition, -1)
	return newString
}
