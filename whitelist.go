package gsoup

import "golang.org/x/net/html/atom"

// Tagset encapsulates a set of unique HTML elements
type Tagset map[atom.Atom]struct{}

// Attrset encapsulates a set of unique element attributes
type Attrset map[string]struct{}

// Protoset encapsulates a set of unique attribute value protocols
type Protoset map[string]struct{}

// Attrmap encapsulates a map of unique element attributes and their values.
// Values will be escaped when set on attributes
type Attrmap map[string]string

// Protomap encapsulates a set of protocols to be enforced on an attribute value
type Protomap map[string]Protoset

// Tagdef encapsulates a single element and its allowed attributes
type Tagdef struct {
	Tag               atom.Atom
	AllowedAttrs      Attrset
	EnforcedAttrs     Attrmap
	EnforcedProtocols Protomap

	// allowRelativeLinks controls whether relative links should be permitted during
	// protocol enforcement. Has no function on attr values where protocols are not
	// enforced via a rule
	allowRelativeLinks bool
}

type whitelist map[atom.Atom]*Tagdef

// EnforceAttr marks an attribute as enforced for a tag. Any tags encountered by the
// parser will have the attribute key and value applied to them.
func (t *Tagdef) EnforceAttr(key string, value string) *Tagdef {
	if t.EnforcedAttrs == nil {
		t.EnforcedAttrs = make(map[string]string)
	}
	t.EnforcedAttrs[normalizeAttrKey(key)] = value
	return t
}

// EnforceProtocols whitelists only the specified protocols for the given attr
// (only applies to the receiver's tag). If protocol enforcement is
// applied, the attr value must be a valid URL per Go's url.Parse() functionality
func (t *Tagdef) EnforceProtocols(attr string, protocols ...string) *Tagdef {
	if t.EnforcedProtocols == nil {
		t.EnforcedProtocols = make(Protomap)
	}
	attr = normalizeAttrKey(attr)
	protoset := make(Protoset)
	for _, proto := range protocols {
		proto = normalizeProtocol(proto)
		if proto != "" {
			protoset[proto] = struct{}{}
		}
	}
	t.EnforcedProtocols[attr] = protoset
	return t
}

// AllowRelativeLinks turns on the ability for links to be relative during protocol
// enforcement. Thus, this setting has no effect unless EnforceProtocols has been called
// on the receiver. Default: false.
func (t *Tagdef) AllowRelativeLinks() *Tagdef {
	t.allowRelativeLinks = true
	return t
}
