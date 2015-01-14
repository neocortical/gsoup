package gsoup

import (
	"strings"

	"golang.org/x/net/html/atom"
)

// Tagset encapsulates a set of unique HTML elements
type Tagset map[atom.Atom]struct{}

// Attrset encapsulates a set of unique element attributes
type Attrset map[string]struct{}

// Attrmap encapsulates a map of unique element attributes and their values
// values will be escaped when set on attributes
type Attrmap map[string]string

// Tagdef encapsulates a single element and its allowed attributes
type Tagdef struct {
	Tag           atom.Atom
	AllowedAttrs  Attrset
	EnforcedAttrs Attrmap
}

type whitelist map[atom.Atom]*Tagdef

// Enforce marks an attribute as enforced for a tag. Any tags encountered by the
// parser will have the attribute key and value applied to them.
func (t *Tagdef) Enforce(key string, value string) *Tagdef {
	if t.EnforcedAttrs == nil {
		t.EnforcedAttrs = make(map[string]string)
	}
	t.EnforcedAttrs[strings.ToLower(key)] = value
	return t
}
