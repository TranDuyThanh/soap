package soap

import (
	"testing"
)

func TestEncodeXML(t *testing.T) {
	xmlString := `<set>
    <longitude>106.825347</longitude>
    <latitude>-6.246246</latitude>
    <serviceTypeId>1</serviceTypeId>
    </set>`

	result := EncodeXML(xmlString)
	expectedString := "&lt;set&gt;&lt;longitude&gt;106.825347&lt;/longitude&gt;&lt;latitude&gt;-6.246246&lt;/latitude&gt;&lt;serviceTypeId&gt;1&lt;/serviceTypeId&gt;&lt;/set&gt;"
	if result != expectedString {
		t.Error("Wrong expected result")
	}
}

func TestDecodeXML(t *testing.T) {
	encodedString := `&lt;set&gt;&lt;longitude&gt;106.825347&lt;/longitude&gt;&lt;latitude&gt;-6.246246&lt;/latitude&gt;&lt;serviceTypeId&gt;1&lt;/serviceTypeId&gt;&lt;/set&gt;`
	expectedString := `<set><longitude>106.825347</longitude><latitude>-6.246246</latitude><serviceTypeId>1</serviceTypeId></set>`

	result := DecodeXML(encodedString)
	if result != expectedString {
		t.Error("Wrong expected result")
	}
}