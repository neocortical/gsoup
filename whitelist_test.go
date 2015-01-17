package gsoup

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html/atom"
)

func Test_EnforceAttr(t *testing.T) {
	tdef := &Tagdef{Tag: atom.P}

	tdef2 := tdef.EnforceAttr("foo", "bar")
	assert.Equal(t, tdef, tdef2, "returned value from factory patter should equal receiver")
	assert.Equal(t, "bar", tdef.EnforcedAttrs["foo"], "tagdef should force 'foo' to 'bar'")

	tdef = tdef.EnforceAttr("BAZ", "qux")
	assert.Equal(t, "qux", tdef.EnforcedAttrs["baz"], "tagdef should lowercase attr keys")
}

func Test_EnforceAttr_Invalid(t *testing.T) {
	tdef := &Tagdef{Tag: atom.Div}

	tdef.EnforceAttr("!@#$a%^&*(b)_+-=|}{\\][\":';?>c<,./~`'\"d}", "value")
	_, ok := tdef.EnforcedAttrs["!@#$a%^&*(b)_+-|}{\\][:;?c<,.~`d}"]
	assert.True(t, ok, "attribute key should be normalized to 'abcd'")
}

func Test_EnforceProtocols(t *testing.T) {
	tdef := &Tagdef{}
	tdef.EnforceProtocols("HREF", "http", "HTTPS", "9chan", "", "+1", "A", "foo:")

	assert.Equal(t, 4, len(tdef.EnforcedProtocols["href"]))
	_, ok := tdef.EnforcedProtocols["href"]["http"]
	assert.True(t, ok, "http should be allowed")
	_, ok = tdef.EnforcedProtocols["href"]["https"]
	assert.True(t, ok, "https should be allowed")
	_, ok = tdef.EnforcedProtocols["href"]["a"]
	assert.True(t, ok, "a should be allowed")
	_, ok = tdef.EnforcedProtocols["href"]["foo"]
	assert.True(t, ok, "foo should be allowed")
}

func Test_EnforceProtocols_overwriteExisting(t *testing.T) {
	tdef := T(atom.A, "href").EnforceProtocols("href", "http", "https")
	tdef2 := tdef.EnforceProtocols("href", "ftp")

	assert.True(t, reflect.DeepEqual(tdef, tdef2))
	assert.True(t, reflect.DeepEqual(tdef.EnforcedProtocols, Protomap{"href": Protoset{"ftp": struct{}{}}}))
}

func Test_AllowRelativeLinks(t *testing.T) {
	tdef := T(atom.P, "class")
	assert.False(t, tdef.allowRelativeLinks)

	tdef.AllowRelativeLinks()
	assert.True(t, tdef.allowRelativeLinks)
}
