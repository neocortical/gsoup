package gsoup

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// XNode is a wrapper around an *html.Node to enforce safe transformations
type XNode interface {
	FirstChild() XNode
	LastChild() XNode
	Type() html.NodeType
	Atom() atom.Atom
	Data() string
	Attr() []html.Attribute
	GetAttr(key string) *html.Attribute
	SetType(html.NodeType)
	SetAtom(atom.Atom)
	SetData(string)
	SetAttrs([]html.Attribute)
}

// TransformFunc describes the signature of a transform function.
type TransformFunc func(XNode) XNode

// PostCleanCallback is called after transformers run and node is cleaned.
// If the node will be deleted based on the whitelist, the second argument will be true.
type PostCleanCallback func(XNode, bool)

func newXNode(n *html.Node) XNode {
	return &tnode{node: n}
}

type tnode struct {
	node *html.Node
}

func (t *tnode) FirstChild() XNode {
	if t.node.FirstChild != nil {
		return newXNode(t.node.FirstChild)
	}
	return nil
}

func (t *tnode) LastChild() XNode {
	if t.node.LastChild != nil {
		return newXNode(t.node.LastChild)
	}
	return nil
}

func (t *tnode) Type() html.NodeType {
	return t.node.Type
}

func (t *tnode) Atom() atom.Atom {
	return t.node.DataAtom
}

func (t *tnode) Data() string {
	return t.node.Data
}

func (t *tnode) Attr() []html.Attribute {
	result := make([]html.Attribute, len(t.node.Attr), len(t.node.Attr))
	copy(result, t.node.Attr)
	return result
}

func (t *tnode) GetAttr(key string) *html.Attribute {
	for _, attr := range t.node.Attr {
		if attr.Key == key {
			var result = html.Attribute{
				Key:       attr.Key,
				Namespace: attr.Namespace,
				Val:       attr.Val,
			}
			return &result
		}
	}

	return nil
}

func (t *tnode) SetType(newType html.NodeType) {
	t.node.Type = newType
}

func (t *tnode) SetAtom(newAtom atom.Atom) {
	t.node.DataAtom = newAtom
	if t.node.Type == html.ElementNode {
		t.node.Data = newAtom.String()
	}
}

func (t *tnode) SetData(newData string) {
	t.node.Data = newData
}

func (t *tnode) SetAttrs(newAttrs []html.Attribute) {
	t.node.Attr = newAttrs
}

func (t *tnode) SetAttr(key, val string) {
	var result = html.Attribute{
		Key:       key,
		Namespace: "",
		Val:       val,
	}

	for i, attr := range t.node.Attr {
		if attr.Key == key {
			t.node.Attr[i] = result
			return
		}
	}

	t.node.Attr = append(t.node.Attr, result)
	return
}
