package gsoup

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Test_cloneWhitelist(t *testing.T) {
	w1 := whitelist{
		atom.H1:  T(atom.H1, "class"),
		atom.Div: T(atom.Div, "id", "class"),
		atom.A:   T(atom.A).EnforceAttr("rel", "nofollow"),
		atom.Img: T(atom.Img, "src").EnforceProtocols("src", "http"),
	}

	w2 := cloneWhitelist(w1)
	assert.True(t, reflect.DeepEqual(w1, w2))

	// manipulate w2
	w2[atom.H2] = T(atom.H2, "onclick")
	assert.Equal(t, 5, len(w2), "w2 should now have 5 tag defs")
	assert.Equal(t, 4, len(w1), "w1 should still have 4 tag defs")

	delete(w2[atom.Div].AllowedAttrs, "id")
	_, present := w1[atom.Div].AllowedAttrs["id"]
	assert.True(t, present, "w1's tagdef should not be mutable thought w2")
}

func Test_normalizeAttrKey(t *testing.T) {
	dirtyKey := "  key/\r\n\t >\"'=name\u0000バナナ \t"
	key := normalizeAttrKey(dirtyKey)

	assert.Equal(t, "keynameバナナ", key, "expected 'keynameバナナ' but got %s", key)
}

func Test_normalizeProtocol(t *testing.T) {
	for raw, expected := range protocolTests {
		actual := normalizeProtocol(raw)
		assert.Equal(t, expected, actual)
	}
}

func TestHTML(t *testing.T) {
	n := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.P,
		Data:     "p",
	}
	nodes, err := html.ParseFragment(strings.NewReader(`This is <em>SO MUCH</em> fun!`), n)
	assert.Nil(t, err)

	for i, n1 := range nodes {
		n1.Parent = n
		if i > 0 {
			n1.PrevSibling = nodes[i-1]
		}
		if i < len(nodes)-1 {
			n1.NextSibling = nodes[i+1]
		}
	}

	n.FirstChild = nodes[0]
	n.LastChild = nodes[len(nodes)-1]

	var buf = &bytes.Buffer{}
	err = html.Render(buf, n)
	fmt.Println(buf.String())

	xnode := newXNode(n)

	actual, err := HTML(xnode)
	assert.Nil(t, err)
	assert.Equal(t, `This is <em>SO MUCH</em> fun!`, actual)
}

var protocolTests = map[string]string{
	"HTTP":              "http",
	"Https":             "https",
	"jav&#x09;ascript:": "javx09ascript",
	"":                  "",
	"###":               "",
	"8chan":             "",
	"+http":             "",
	"a":                 "a",
	"A":                 "a",
	"Z":                 "z",
}
