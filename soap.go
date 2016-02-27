package soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type SoapRequest struct {
	UrlString string
	Template  string
	Action    string
	Condition interface{}
}

var code = map[string]string{
	"&":  "&amp;",
	"<":  "&lt;",
	">":  "&gt;",
	"'":  "&apos;",
	"\"": "&quot;",
}

func (this *SoapRequest) SendAndParseResponseTo(v interface{}) {
	encodedCondition := createConditionString(this.Condition)
	payload := initPayload(this.Template, this.Action, encodedCondition)
	payloadByte := bytes.NewBuffer([]byte(payload))
	req, err := http.NewRequest("POST", this.UrlString, payloadByte)
	check(err)
	req.Header.Add("Content-Type", "text/xml")
	client := http.Client{}
	resp, err := client.Do(req)
	check(err)
	if resp.StatusCode != 200 {
		fmt.Errorf("Response status is %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	resultInfo := getUnescapedResponse(string(body))
	xml.Unmarshal([]byte(resultInfo), v)
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

func createConditionString(condition interface{}) string {
	conditionStr := parseToXml(condition)
	encodedCondition := EscapeXML(conditionStr)
	return encodedCondition
}

func getResponse(xmlString string) string {
	r, _ := regexp.Compile("(&lt;resultInfo)(.+)(resultInfo&gt;)")
	newString := r.FindString(xmlString)
	return newString
}

func getUnescapedResponse(xmlString string) string {
	xmlString = getResponse(xmlString)
	newString := UnescapeXML(xmlString)
	return newString
}

func initPayload(template, action, condition string) string {
	newString := strings.Replace(template, "*ACTION*", action, -1)
	newString = strings.Replace(newString, "*CONDITION*", condition, -1)
	return newString
}

func parseToXml(v interface{}) string {
	conditionXml, err := xml.Marshal(&v)
	check(err)
	return string(conditionXml)
}

func check(err error) {
	if err != nil {
		fmt.Errorf("%v", err)
	}
}
