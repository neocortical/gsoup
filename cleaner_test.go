package gsoup

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Test_FailGracefullyOnParseError(t *testing.T) {
	_, err := NewEmptyCleaner().Clean(badReader{})
	assert.NotNil(t, err, "err should not be nil")
}

func Test_stripInvalidAttributes(t *testing.T) {
	// basic passthrough
	elem := ele("class")
	def := T(atom.P, "class")
	stripInvalidAttributes(elem, def)
	assert.Equal(t, 1, len(elem.Attr), "tag should have one attribute")
	assert.Equal(t, "class", elem.Attr[0].Key, "elem should still contain key 'class'")

	// basic strip
	def = T(atom.P)
	stripInvalidAttributes(elem, def)
	assert.Equal(t, 0, len(elem.Attr), "tag should have zero attributes")

	// attributes should be found case insensitive and lowercased
	def = T(atom.P, "class")
	elem = ele("ClAsS", "OnClicK")
	stripInvalidAttributes(elem, def)
	assert.Equal(t, 1, len(elem.Attr), "tag should have one attribute")
	assert.Equal(t, "class", elem.Attr[0].Key, "elem should contain lowercased key 'class'")
}

func Test_shouldPreserveChildren(t *testing.T) {
	c := NewBasicCleaner().(*cleaner)

	// test 'always preserve' elements
	n := ele()
	n.Type = html.ElementNode
	for a := range preserveChildrenSet {
		n.DataAtom = a
		assert.True(t, c.shouldPreserveChildren(n), "should preserve atoms in the 'always preserve set'")
	}

	// test don't preserve children of non-element types
	n.Type = html.TextNode
	assert.False(t, c.shouldPreserveChildren(n), "should not preserve non-elements")

	// test preserveChildren toggle
	n.Type = html.ElementNode
	n.DataAtom = atom.P
	c.preserveChildren = false
	assert.False(t, c.shouldPreserveChildren(n), "should not preserve P tag if preserveChildren is false")
	c.preserveChildren = true
	assert.True(t, c.shouldPreserveChildren(n), "should preserve P tag if preserveChildren is true")

	// always delete non-structure nodes
	for a := range deleteChildrenSet {
		n.DataAtom = a
		assert.False(t, c.shouldPreserveChildren(n), "should not preserve atoms in the 'always delete set'")
	}
}

func Test_removeElement_WithChildren(t *testing.T) {
	p := &html.Node{Type: html.DocumentNode}
	n1 := eleWithData(1)
	n2 := eleWithData(2)
	n3 := eleWithData(3)
	c1 := eleWithData(4)
	c2 := eleWithData(5)
	p.AppendChild(n1)
	p.AppendChild(n2)
	p.AppendChild(n3)
	n2.AppendChild(c1)
	n2.AppendChild(c2)
	n2.Type = html.ElementNode
	n2.DataAtom = atom.P

	c := NewBasicCleaner().(*cleaner)
	c.preserveChildren = true

	result := c.removeElement(n2)
	assert.Equal(t, result, c1, "returned node should be first child of removed node if it has children")
	assert.Nil(t, n2.Parent, "n2 should have been removed from the graph")
	assert.Equal(t, n1, c1.PrevSibling, "c1 should be next sib of n1")
	assert.Equal(t, c2, c1.NextSibling, "c2 should be next sib of c1")
	assert.Equal(t, n3, c2.NextSibling, "n3 should be next sib of c2")
}

func Test_removeElement_WithoutChildren(t *testing.T) {
	p := &html.Node{Type: html.DocumentNode}
	n1 := eleWithData(1)
	n2 := eleWithData(2)
	n3 := eleWithData(3)
	p.AppendChild(n1)
	p.AppendChild(n2)
	p.AppendChild(n3)
	n2.Type = html.ElementNode
	n2.DataAtom = atom.P

	c := NewBasicCleaner().(*cleaner)
	c.preserveChildren = true

	result := c.removeElement(n2)
	assert.Nil(t, n2.Parent, "n2 should have been removed from the graph")
	assert.Equal(t, n3, result, "returned node should be next sib. of removed node if removed node has no children")
	assert.Equal(t, n3, n1.NextSibling, "n1 and n3 should now be neighbors")
}

func Test_removeElement_AndKillChildren(t *testing.T) {
	p := &html.Node{Type: html.DocumentNode}
	n1 := eleWithData(1)
	n2 := eleWithData(2)
	n2.DataAtom = atom.Script // don't preserve children of this element
	n3 := eleWithData(3)
	c1 := eleWithData(4)
	c2 := eleWithData(5)
	p.AppendChild(n1)
	p.AppendChild(n2)
	p.AppendChild(n3)
	n2.AppendChild(c1)
	n2.AppendChild(c2)

	c := NewBasicCleaner().(*cleaner)

	result := c.removeElement(n2)
	assert.Equal(t, result, n3, "returned node should be next sibling of removed node")
	assert.Nil(t, n2.Parent, "n2 should have been removed from the graph")
	assert.Equal(t, n1, n3.PrevSibling, "n3 should be next sib of n1")
}

func Test_removeElement_ReturnsNilIfItWasLastChild(t *testing.T) {
	p := &html.Node{Type: html.DocumentNode}
	n1 := eleWithData(1)
	n2 := eleWithData(2)
	p.AppendChild(n1)
	p.AppendChild(n2)

	c := NewBasicCleaner().(*cleaner)

	result := c.removeElement(n2)
	assert.Nil(t, n2.Parent, "n2 should have been removed from the graph")
	assert.Nil(t, result, "returned node should be nil if it was a last child")
	assert.Nil(t, n1.NextSibling, "n1 should have no more siblings")
}

func Test_Clean_All(t *testing.T) {
	c := NewBasicCleaner().(*cleaner)

	for input, expected := range basicWhitelistKillChildren {
		doc, err := c.Clean(strings.NewReader(input))
		var buf bytes.Buffer
		html.Render(&buf, doc)
		actual := buf.String()
		assert.Nil(t, err, "unexpected error: %v", err)
		assert.Equal(t, expected, actual, "expected %s but got %s", expected, actual)
	}

	c.preserveChildren = true
	for input, expected := range basicWhitelistpreserveChildren {
		doc, err := c.Clean(strings.NewReader(input))
		var buf bytes.Buffer
		html.Render(&buf, doc)
		actual := buf.String()
		assert.Nil(t, err, "unexpected error: %v", err)
		assert.Equal(t, expected, actual, "expected %s but got %s", expected, actual)
	}
}

func Test_AddTags(t *testing.T) {
	c := NewBasicCleaner().(*cleaner)

	c.AddTags(
		T(atom.Div, "id"),
		T(atom.Table),
		T(atom.Q, "id"),
	)

	def, ok := c.w[atom.Div]
	assert.True(t, ok, "div tag should now appear in whitelist")
	_, ok = def.AllowedAttrs["id"]
	assert.True(t, ok, "attribute 'id' should be set for div tag")

	def, ok = c.w[atom.Table]
	assert.True(t, ok, "table tag should now appear in whitelist")

	def, ok = c.w[atom.Q]
	assert.True(t, ok, "q tag should still appear in whitelist")
	_, ok = def.AllowedAttrs["id"]
	assert.True(t, ok, "attribute 'id' should be set for q tag")
	_, ok = def.AllowedAttrs["cite"]
	assert.False(t, ok, "attribute 'cite' should no longer be set for q tag")
}

func Test_RemoveTags(t *testing.T) {
	c := NewBasicCleaner().(*cleaner)

	c.RemoveTags(atom.P, atom.Div)

	_, ok := c.w[atom.P]
	assert.False(t, ok, "p tag should no longer appear in whitelist")
	_, ok = c.w[atom.Div]
	assert.False(t, ok, "div tag should still not appear in whitelist")
}

func Test_PreserveChildren(t *testing.T) {
	c := &cleaner{}
	assert.False(t, c.preserveChildren, "default should be false")
	c2 := c.PreserveChildren()
	assert.True(t, c.preserveChildren, "default should be false")
	assert.True(t, c2.(*cleaner).preserveChildren, "default should be false for returned value")
}

func Test_Clean_ShouldOverwriteEnforcedAttribute(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.A, "rel").Enforce("rel", "nofollow"))
	input := `<a rel="foobar">hello</a>`
	doc, err := c.Clean(strings.NewReader(input))
	assert.Nil(t, err, "err should be nil")
	var buf bytes.Buffer
	html.Render(&buf, doc)
	actual := buf.String()
	assert.Equal(t, `<a rel="nofollow">hello</a>`, actual, "should overwrite enforced attributes")
}

func ele(attrs ...string) *html.Node {
	attributes := []html.Attribute{}
	for _, key := range attrs {
		attributes = append(attributes, html.Attribute{Key: key})
	}
	return &html.Node{
		Attr: attributes,
	}
}

func eleWithData(datum int) *html.Node {
	return &html.Node{
		Data: strconv.Itoa(datum),
	}
}

var basicWhitelistpreserveChildren = map[string]string{
	`plain text`:                                        `plain text`,
	`plain text<!-- comment -->`:                        `plain text`,
	`<p>plain text</p><div>more text</div>`:             `<p>plain text</p>more text`,
	`<SCRIPT SRC=http://ha.ckers.org/xss.js></SCRIPT>`:  ``,
	`<IMG SRC="javascript:alert('XSS');">`:              ``,
	`<IMG """><SCRIPT>alert("XSS")</SCRIPT>">`:          `&#34;&gt;`,
	`<P """><SCRIPT>alert("XSS")</SCRIPT>">`:            `<p>&#34;&gt;</p>`,
	`<p onmouseover="alert('xxs')">`:                    `<p></p>`,
	`<P/XSS SRC="http://ha.ckers.org/xss.js"></SCRIPT>`: `<p></p>`,
	`<<SCRIPT>alert("XSS");//<</SCRIPT>`:                `&lt;`,
	`<BR SIZE="&{alert('XSS')}">`:                       `<br/>`,
	`exp/*<A STYLE='no\xss:noxss("*//*");
	xss:ex/*XSS*//*/*/pression(alert("XSS"))'>`: `exp/*<a rel="nofollow"></a>`,
	`<!--[if gte IE 4]>
	<SCRIPT>alert('XSS');</SCRIPT>
	<![endif]-->`: ``,
	`<a onmouseover="alert(document.cookie)">xxs link</a>`: `<a rel="nofollow">xxs link</a>`,
	`<a onmouseover=alert(document.cookie)>xxs link</a>`:   `<a rel="nofollow">xxs link</a>`,
	`<a rel="foobar">http://google.com</a>`:                `<a rel="nofollow">http://google.com</a>`,
}

var basicWhitelistKillChildren = map[string]string{
	`plain text`:                            `plain text`,
	`plain text<!-- comment -->`:            `plain text`,
	`<p>plain text</p><div>more text</div>`: `<p>plain text</p>`,
}

type badReader struct{}

func (br badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("i've made a terrible mistake")
}
