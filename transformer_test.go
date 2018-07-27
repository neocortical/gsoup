package gsoup

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Test_ShouldNotAllowTransformingDocNode(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.I))
	c.AddTransformer(func(x XNode) XNode {
		if x.Type() == html.DocumentNode {
			x.SetType(html.TextNode)
			x.SetData("wot")
		}
		return x
	})
	actual, err := c.CleanString(`<i>some text</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>some text</i>`, actual)
}

func Test_ShouldDoSimpleTagReplace(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.B))
	c.AddTransformer(func(x XNode) XNode {
		if x.Atom() == atom.I {
			x.SetAtom(atom.B)
		}
		return x
	})
	actual, err := c.CleanString(`<i>some text</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `<b>some text</b>`, actual)
}

func Test_ShouldStillValidateTransformedElements(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.B)).PreserveChildren()
	c.AddTransformer(func(x XNode) XNode {
		if x.Atom() == atom.I {
			x.SetAtom(atom.Em)
		}
		return x
	})
	actual, err := c.CleanString(`<i>some text</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `some text`, actual)
}

func Test_ShouldReplaceText(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.B))
	c.AddTransformer(func(x XNode) XNode {
		if x.Type() == html.TextNode {
			x.SetData("xxx")
		}
		return x
	})
	actual, err := c.CleanString(`<b>some text</b>`)
	assert.Nil(t, err)
	assert.Equal(t, `<b>xxx</b>`, actual)
}

func Test_ShouldDeleteNodeAndChildren(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.B), T(atom.I))
	c.AddTransformer(func(x XNode) XNode {
		if x.Atom() == atom.B {
			return nil
		}
		return x
	})
	actual, err := c.CleanString(`<i>keep</i><b><i>destroy</i></b><i>keep2</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>keep</i><i>keep2</i>`, actual)
}

func Test_ShouldAllowReplacementWithFirstChild(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.B), T(atom.I))
	c.AddTransformer(func(x XNode) XNode {
		if x.Atom() == atom.B {
			return x.FirstChild()
		}
		return x
	})

	actual, err := c.CleanString(`<i>keep1</i><b><i>keep2</i><i>forget about me</i></b><i>keep3</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>keep1</i><i>keep2</i><i>keep3</i>`, actual)

	// test that it works for a nil child
	actual, err = c.CleanString(`<i>keep1</i><b></b><i>keep3</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>keep1</i><i>keep3</i>`, actual)
}

func Test_ShouldAllowReplacementWithLastChild(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.B), T(atom.I))
	c.AddTransformer(func(x XNode) XNode {
		if x.Atom() == atom.B {
			return x.LastChild()
		}
		return x
	})

	actual, err := c.CleanString(`<i>keep1</i><b><i>forget about me</i>me too<i>keep2</i></b><i>keep3</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>keep1</i><i>keep2</i><i>keep3</i>`, actual)

	// test that it works for a nil child
	actual, err = c.CleanString(`<i>keep1</i><b></b><i>keep3</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>keep1</i><i>keep3</i>`, actual)
}

func Test_ShouldInspectDataAndConditionallyReplace(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.I))
	c.AddTransformer(func(x XNode) XNode {
		if x.Type() == html.TextNode {
			x.SetData(strings.ToUpper(x.Data()))
		}
		return x
	})
	actual, err := c.CleanString(`<i>one</i>two<i>three</i>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>ONE</i>TWO<i>THREE</i>`, actual)
}

func Test_ShouldBeAbleToModifyElementAttributes(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.I, "foo", "bar"), T(atom.B, "foo", "bar"))
	c.AddTransformer(func(x XNode) XNode {
		if x.Atom() == atom.I {
			attrs := x.Attr()
			var newAttrs []html.Attribute
			for _, attr := range attrs {
				if attr.Key == "foo" {
					newAttrs = append(newAttrs, html.Attribute{Key: "foo", Val: "manchu"})
				}
			}
			x.SetAttrs(newAttrs)
		}
		return x
	})
	actual, err := c.CleanString(`<i foo="fighter" bar="brawl">hi</i><b foo="fighter" bar="brawl">there</b>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i foo="manchu">hi</i><b foo="fighter" bar="brawl">there</b>`, actual)

}

func Test_ShouldBeAbleToSetType(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.B), T(atom.I))
	c.AddTransformer(func(x XNode) XNode {
		if x.Type() == html.ElementNode && x.Atom() == atom.B {
			x.SetType(html.TextNode)
		}
		return x
	})
	actual, err := c.CleanString(`<i>a</i><b>c</b>d`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>a</i>bd`, actual)
}

func Test_ShouldApplyMultipleTransformers(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.B), T(atom.I))
	c.AddTransformer(func(x XNode) XNode {
		if x.Type() == html.ElementNode && x.Atom() == atom.Em {
			x.SetAtom(atom.I)
		}
		return x
	})
	c.AddTransformer(func(x XNode) XNode {
		if x.Type() == html.ElementNode && x.Atom() == atom.Strong {
			x.SetAtom(atom.B)
		}
		return x
	})
	actual, err := c.CleanString(`<em>a<strong>b</strong></em><strong>c</strong>`)
	assert.Nil(t, err)
	assert.Equal(t, `<i>a<b>b</b></i><b>c</b>`, actual)

}

func Test_EmptyElement(t *testing.T) {
	c := NewEmptyCleaner().AddTags(T(atom.Div), T(atom.I))
	c.AddTransformer(func(x XNode) XNode {
		if x.Type() == html.ElementNode && x.Atom() == atom.Div {
			assert.Nil(t, x.FirstChild())
		}
		return x
	})
	actual, err := c.CleanString(`<div></div>`)
	assert.Nil(t, err)
	assert.Equal(t, `<div></div>`, actual)
}

func TestGetAttr(t *testing.T) {
	node := &tnode{
		node: &html.Node{
			Attr: []html.Attribute{
				html.Attribute{Key: "foo", Namespace: "bar", Val: "baz"},
				html.Attribute{Key: "a", Namespace: "b", Val: "c"},
			},
		},
	}

	attr := node.GetAttr("a")
	assert.NotNil(t, attr)
	assert.Equal(t, "a", attr.Key)
	assert.Equal(t, "b", attr.Namespace)
	assert.Equal(t, "c", attr.Val)

	// object should be cloned
	attr.Val = "f"
	assert.Equal(t, "c", node.node.Attr[1].Val)

	attr = node.GetAttr("badkey")
	assert.Nil(t, attr)
}
