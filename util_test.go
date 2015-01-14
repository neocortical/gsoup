package gsoup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html/atom"
)

func Test_cloneWhitelist(t *testing.T) {
	w1 := whitelist{
		atom.H1:  T(atom.H1, "class"),
		atom.Div: T(atom.Div, "id", "class"),
	}

	w2 := cloneWhitelist(w1)
	assert.Equal(t, 2, len(w2), "w2 should have 2 tag defs")
	assert.Equal(t, w1[atom.H1], w2[atom.H1], "w1 and w2 should have equal values")
	assert.Equal(t, w1[atom.Div], w2[atom.Div], "w1 and w2 should have equal values")

	// manipulate w2
	w2[atom.H2] = T(atom.H2, "onclick")
	assert.Equal(t, 3, len(w2), "w2 should now have 3 tag defs")
	assert.Equal(t, 2, len(w1), "w1 should still have 2 tag defs")

	delete(w2[atom.Div].AllowedAttrs, "id")
	_, present := w1[atom.Div].AllowedAttrs["id"]
	assert.True(t, present, "w1's tagdef should not be mutable thought w2")
}