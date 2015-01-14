package gsoup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html/atom"
)

func Test_T(t *testing.T) {
	// test that attributes are lower-cased
	def := T(atom.Div, "ID", "ClAsS", "foo")

	assert.True(t, len(def.AllowedAttrs) == 3, "tagdef should have 3 attributes")
	_, ok := def.AllowedAttrs["id"]
	assert.True(t, ok, "'id' should be an attribute")
	_, ok = def.AllowedAttrs["class"]
	assert.True(t, ok, "'class' should be an attribute")
	_, ok = def.AllowedAttrs["foo"]
	assert.True(t, ok, "'foo' should be an attribute")
	_, ok = def.AllowedAttrs["bar"]
	assert.True(t, !ok, "'bar' should not be an attribute")
	_, ok = def.AllowedAttrs["ID"]
	assert.True(t, !ok, "'ID' should not be an attribute")
}

func Test_T_EscapesAttrKeys(t *testing.T) {
	def := T(atom.Div, "  key/\r\n\t >\"'=name\u0000バナナ \t")
	_, normalizedKeyExists := def.AllowedAttrs["keynameバナナ"]
	assert.Equal(t, 1, len(def.AllowedAttrs), "tagdef should only have one allowed attr")
	assert.True(t, normalizedKeyExists, "expected 'keynameバナナ' but key doesn't exist")
}
