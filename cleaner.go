package gsoup

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Cleaner defines the interface for sanitizing markup
type Cleaner interface {
	// Clean sanizitizes HTML input based on the cleaner's rules
	Clean(string) (*html.Node, error)
	// AddTags adds acceptable tags (and their allowed attributes) to the whitelist
	AddTags(tags ...Tagdef) Cleaner
	// RemoveTags removes tags that should be deleted during sanitization
	RemoveTags(tags ...atom.Atom) Cleaner
	// PreserveChildren causes child nodes of deleted tags to be retained (if they themselves are allowed)
	PreserveChildren() Cleaner
}

type cleaner struct {
	w whitelist

	// preserveChildren controls whether children of deleted nodes are also deleted. This
	// setting does not apply to elements that can contain no user-facing text (e.g. <script>)
	// Default: false
	preserveChildren bool
}

func (c *cleaner) Clean(input string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return doc, err
	}

	c.cleanRecursive(doc)

	return doc, nil
}

func (c *cleaner) AddTags(tags ...Tagdef) Cleaner {
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

// cleanRecursive performs a depth-first traversal of the DOM, removing nodes and attributes in place as it goes
func (c *cleaner) cleanRecursive(n *html.Node) *html.Node {
	switch n.Type {
	case html.ElementNode:
		tagdef, ok := c.w[n.DataAtom]
		if !ok {
			return c.removeElement(n)
		}

		stripInvalidAttributes(n, &tagdef)

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
	newAttr := n.Attr[:0]
	for _, attr := range n.Attr {
		_, attrAllowed := tagdef.AllowedAttrs[strings.ToLower(attr.Key)]
		if attrAllowed {
			attr.Key = strings.ToLower(attr.Key)
			newAttr = append(newAttr, attr)
		}
	}
	n.Attr = newAttr
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
