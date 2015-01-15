package gsoup

import (
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
