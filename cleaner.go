package gsoup

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Cleaner defines the interface for sanitizing markup
type Cleaner interface {
	// Clean sanizitizes HTML input based on the cleaner's rules
	Clean(io.Reader) (*html.Node, error)
	// CleanNode sanitizes HTML in an already constructed document. NOTE: This call is
	// destructive to the input param. The returned node will be wrapped in a
	// DocumentNode if the input wasn't already.
	CleanNode(root *html.Node) (*html.Node, error)
	// CleanString is a convenience wrapper for simple, string-in-string-out cleaning of markup
	CleanString(input string) (string, error)
	// AddTags adds acceptable tags (and their allowed attributes) to the whitelist
	AddTags(tags ...*Tagdef) Cleaner
	// RemoveTags removes tags that should be deleted during sanitization
	RemoveTags(tags ...atom.Atom) Cleaner
	// PreserveChildren causes child nodes of deleted tags to be retained (if they themselves are allowed)
	PreserveChildren() Cleaner

	AddTransformer(TransformFunc) Cleaner
}

type cleaner struct {
	// the whitelist of allowed tags and their allowed attributes
	w whitelist

	// preserveChildren controls whether children of deleted nodes are also deleted. This
	// setting does not apply to elements that can contain no user-facing text (e.g. <script>)
	// Default: false
	preserveChildren bool

	// transforms is a list of transforms registered with this cleaner
	transforms []TransformFunc
}

var errorInvalidProtocol = errors.New("invalid protocol")
var errorRelativeLink = errors.New("relative links disallowed")

func (c *cleaner) Clean(input io.Reader) (*html.Node, error) {
	doc, err := html.Parse(input)
	if err != nil {
		return doc, err
	}

	c.cleanRecursive(doc)

	return doc, nil
}

func (c *cleaner) CleanNode(root *html.Node) (*html.Node, error) {
	if root == nil {
		return root, errors.New("root cannot be nil")
	}
	var doc *html.Node
	if root.Type != html.DocumentNode {
		doc = &html.Node{Type: html.DocumentNode}
		doc.AppendChild(root)
	} else {
		doc = root
	}

	c.cleanRecursive(doc)

	return doc, nil
}

func (c *cleaner) CleanString(input string) (string, error) {
	doc, err := c.Clean(strings.NewReader(input))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *cleaner) AddTags(tags ...*Tagdef) Cleaner {
	for _, tagdef := range tags {
		c.w[tagdef.Tag] = tagdef
	}
	return c
}

func (c *cleaner) RemoveTags(tags ...atom.Atom) Cleaner {
	for _, tag := range tags {
		delete(c.w, tag)
	}
	return c
}

func (c *cleaner) PreserveChildren() Cleaner {
	c.preserveChildren = true
	return c
}

func (c *cleaner) AddTransformer(t TransformFunc) Cleaner {
	c.transforms = append(c.transforms, t)
	return c
}

// cleanRecursive performs a depth-first traversal of the DOM, removing nodes and attributes in place as it goes
func (c *cleaner) cleanRecursive(n *html.Node) *html.Node {

	// apply any transform functions
	if n.Type == html.ElementNode || n.Type == html.TextNode {
		for _, transform := range c.transforms {
			transformed := transform(newXNode(n))
			if transformed == nil {
				n.Parent.RemoveChild(n)
			} else {
				newNode := transformed.(*tnode).node
				if newNode != n {
					if newNode.Parent != nil {
						newNode.Parent.RemoveChild(newNode)
					}
					n.Parent.InsertBefore(newNode, n)
					n.Parent.RemoveChild(n)
					n = newNode
				}
			}
		}
	}

	switch n.Type {
	case html.ElementNode:
		tagdef, ok := c.w[n.DataAtom]
		if !ok {
			return c.removeElement(n)
		}

		stripInvalidAttributes(n, tagdef)

	case html.ErrorNode, html.CommentNode, html.DoctypeNode:
		return c.removeElement(n)
	}

	ch := n.FirstChild
	for ch != nil {
		ch = c.cleanRecursive(ch)
	}

	return n.NextSibling
}

// stripInvalidAttributes removes non-whitelisted attributes on the node in place
func stripInvalidAttributes(n *html.Node, tagdef *Tagdef) {
	attrMap := make(map[string]int)
	newAttr := n.Attr[:0]
	for _, attr := range n.Attr {
		normalizedAttr := normalizeAttrKey(attr.Key)
		_, attrAllowed := tagdef.AllowedAttrs[normalizedAttr]
		if attrAllowed {
			normalizedVal, err := enforceProtocol(tagdef, normalizedAttr, attr.Val)
			if err == nil {
				attr.Key = normalizedAttr
				attr.Val = normalizedVal
				newAttr = append(newAttr, attr)
				attrMap[attr.Key] = len(newAttr) - 1
			}
		}
	}

	// add any enforced attributes
	for key, value := range tagdef.EnforcedAttrs {
		index, ok := attrMap[key]
		if ok {
			newAttr[index].Val = value
		} else {
			attr := html.Attribute{Key: key, Val: value}
			newAttr = append(newAttr, attr)
		}
	}

	n.Attr = newAttr
}

func enforceProtocol(tagdef *Tagdef, attrKey string, attrVal string) (string, error) {
	if tagdef.EnforcedProtocols == nil {
		return attrVal, nil
	}
	allowedProtos, enforce := tagdef.EnforcedProtocols[attrKey]
	if !enforce {
		return attrVal, nil
	}

	// url must be parsable to be valid
	u, err := url.Parse(attrVal)
	if err != nil {
		return "", err
	}

	// relative link logic
	if !u.IsAbs() && !tagdef.allowRelativeLinks {
		return "", errorRelativeLink
	}
	if !u.IsAbs() {
		return u.String(), nil
	}

	for allowedProto := range allowedProtos {
		if u.Scheme == allowedProto {
			return u.String(), nil
		}
	}

	return "", errorInvalidProtocol
}

func (c *cleaner) removeElement(n *html.Node) (result *html.Node) {
	p := n.Parent

	preserveChildren := c.shouldPreserveChildren(n)
	if preserveChildren {
		result = n.FirstChild
	}

	if result == nil {
		result = n.NextSibling
	}

	for preserveChildren && n.FirstChild != nil {
		ch := n.FirstChild
		n.RemoveChild(ch)
		p.InsertBefore(ch, n)
	}
	p.RemoveChild(n)

	return result
}

func (c *cleaner) shouldPreserveChildren(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}

	_, alwaysPreserve := preserveChildrenSet[n.DataAtom]
	if alwaysPreserve {
		return true
	}

	if !c.preserveChildren {
		return false
	}

	_, antiresult := deleteChildrenSet[n.DataAtom]
	return !antiresult
}
