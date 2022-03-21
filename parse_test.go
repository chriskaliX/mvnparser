package mvnparse

import (
    "encoding/xml"
    "github.com/elliotchance/orderedmap"
    "github.com/stretchr/testify/assert"
    "github.com/subchen/go-xmldom"
    "html"
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
    assert.NoError(t, err)
    result := string(data)
    assert.True(t, strings.Index(result, "b") < strings.Index(result, "a"), "unexpected order for properties")
}

func TestProperties_UnmarshalXML(t *testing.T) {
    data := "<properties><a>1</a><b>1</b><c>1</c></properties>"
    p := Properties{}
    err := xml.Unmarshal([]byte(data), &p)
    assert.NoError(t, err)
    keys := p.Entries.Keys()
    assert.Len(t, keys, 3)
    assert.Equal(t, "a", keys[0])
}

func makeAttribute(name string) *xmldom.Attribute {
    return &xmldom.Attribute{
        Name:  name,
        Value: name,
    }
}

func makeNode(name string) *xmldom.Node {
    return &xmldom.Node{
        Name: name,
        Attributes: []*xmldom.Attribute{
            makeAttribute("attr1"),
            makeAttribute("attr2"),
        },
    }
}

func TestConfiguration_MarshalXML(t *testing.T) {
    children := []*xmldom.Node{
        makeNode("node1"),
        makeNode("node2"),
    }
    children[1].Text = "yyyy-MM-dd'T'HH:mm:ssZ"

    config := Configuration{
        XMLName: xml.Name{
            Space: "",
            Local: "configuration",
        },
        Children: children,
    }
    data, err := xml.Marshal(&config)
    assert.NoError(t, err)
    result := html.UnescapeString(string(data))
    assert.Contains(t, result, `<node1 attr1="attr1"`)
    assert.Contains(t, result, `attr2="attr2"`)
    assert.Contains(t, result, `yyyy-MM-dd'T'HH:mm:ssZ`)
}

func TestConfiguration_UnmarshalXML(t *testing.T) {
    data := `<configuration>
                <transformers>
                  <transformer implementation="org.apache.maven.plugins.shade.resource.AppendingTransformer">
                    <resource>META-INF/spring.handlers</resource>
                  </transformer>
                  <transformer implementation="org.springframework.boot.maven.PropertiesMergingResourceTransformer">
                    <resource>META-INF/spring.factories</resource>
                  </transformer>
                </transformers>
                <test>yyyy-MM-dd'T'HH:mm:ssZ</test>
            </configuration>`
    var config Configuration
    err := xml.Unmarshal([]byte(data), &config)
    assert.NoError(t, err)

    children := config.Children
    assert.NotNil(t, children)
    assert.Len(t, children, 2)

    transforms := children[0]
    assert.Equal(t, "transformers", transforms.Name)
    assert.Len(t, transforms.Children, 2)

    transform := transforms.Children[0]
    assert.Equal(t, "transformer", transform.Name)
    assert.Len(t, transform.Attributes, 1)

    attr := transform.Attributes[0]
    assert.Equal(t, "implementation", attr.Name)
    assert.Equal(t, "org.apache.maven.plugins.shade.resource.AppendingTransformer", attr.Value)

    resources := transform.Children
    assert.Len(t, resources, 1)
    resource := resources[0]
    assert.Equal(t, "resource", resource.Name)
    assert.Equal(t, "META-INF/spring.handlers", resource.Text)

    test := children[1]
    assert.NotNil(t, test)
    assert.Equal(t, "yyyy-MM-dd'T'HH:mm:ssZ", test.Text)
}