package gsoup

import "golang.org/x/net/html/atom"

func cloneWhitelist(in whitelist) (out whitelist) {
	out = make(map[atom.Atom]*Tagdef)
	for tag, tagdef := range in {
		newdef := &Tagdef{Tag: tagdef.Tag, AllowedAttrs: make(Attrset)}
		for attr := range tagdef.AllowedAttrs {
			newdef.AllowedAttrs[attr] = struct{}{}
		}
		for key, value := range tagdef.EnforcedAttrs {
			if newdef.EnforcedAttrs == nil {
				newdef.EnforcedAttrs = make(map[string]string)
			}
			newdef.EnforcedAttrs[key] = value
		}
		out[tag] = newdef
	}
	return out
}
