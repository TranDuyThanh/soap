package soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
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

var patterns = []string{"resultInfo", "orderInfo"}

func (this *SoapRequest) SendAndParseResponseTo(v interface{}) error {
	encodedCondition := createConditionString(this.Condition)

	payload := initPayload(this.Template, this.Action, encodedCondition)
	payloadByte := bytes.NewBuffer([]byte(payload))

	req, err := http.NewRequest("POST", this.UrlString, payloadByte)
	check(err)
	req.Header.Add("Content-Type", "text/xml")

	timeout := time.Duration(30000 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		fmt.Errorf("Response status is %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	resultInfo := getUnescapedResponse(string(body))
	fmt.Println(resultInfo)
	xml.Unmarshal([]byte(resultInfo), v)
	return nil
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

func GetResponse(xmlString string) string {
	for _, pattern := range patterns {
		regexStr := fmt.Sprintf("(&lt;%s)(.+)(%s&gt;)", pattern, pattern)
		r, _ := regexp.Compile(regexStr)
		newString := r.FindString(xmlString)
		if newString != "" {
			return newString
		}
	}

	return ""
}

func getUnescapedResponse(xmlString string) string {
	xmlString = GetResponse(xmlString)
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
		fmt.Printf("%#v", err)
	}
}
