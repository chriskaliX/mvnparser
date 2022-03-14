package mvnparse

import (
	"encoding/xml"
	"github.com/elliotchance/orderedmap"
	"strings"
	"testing"
)

func TestProperties_MarshalXML(t *testing.T) {

	p := Properties{
		Entries: *orderedmap.NewOrderedMap(),
	}
	p.Entries.Set("a", "1")
	p.Entries.Set("b", "1")
	p.Entries.Set("c", "1")
	data, err := xml.Marshal(p)
	if err != nil {
		t.Error("properties marshal failed, ", err)
		return
	}

	result := string(data)
	if strings.Index(result, "b") < strings.Index(result, "a") {
		t.Error("unexpected order for properties, ", result)
		return
	}
}

func TestProperties_UnmarshalXML(t *testing.T) {
	data := "<properties><a>1</a><b>1</b><c>1</c></properties>"
	p := Properties{}
	err := xml.Unmarshal([]byte(data), &p)
	if err != nil {
		t.Error("properties unmarshal failed, ", err)
		return
	}
	keys := p.Entries.Keys()
	if len(keys) != 3 {
		t.Error("mismatch count, exepect 3, got ", len(keys))
		return
	}
	if keys[0] != "a" {
		t.Error("unexpected order for properties", p)
		return
	}
}

