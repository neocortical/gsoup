package gsoup

import "golang.org/x/net/html/atom"

// Tagset encapsulates a set of unique HTML elements
type Tagset map[atom.Atom]struct{}

// Attrset encapsulates a set of unique element attributes
type Attrset map[string]struct{}

// Tagdef encapsulates a single element and its allowed attributes
type Tagdef struct {
	Tag          atom.Atom
	AllowedAttrs Attrset
}

type whitelist map[atom.Atom]Tagdef
