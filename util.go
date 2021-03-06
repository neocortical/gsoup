package gsoup

import (
	"strings"

	"golang.org/x/net/html/atom"
)

func cloneWhitelist(in whitelist) (out whitelist) {
	out = make(map[atom.Atom]*Tagdef)
	for tag, tagdef := range in {
		newdef := &Tagdef{Tag: tagdef.Tag, AllowedAttrs: make(Attrset)}
		for attr := range tagdef.AllowedAttrs {
			newdef.AllowedAttrs[attr] = struct{}{}
		}

		// enforced attrs
		for key, value := range tagdef.EnforcedAttrs {
			if newdef.EnforcedAttrs == nil {
				newdef.EnforcedAttrs = make(map[string]string)
			}
			newdef.EnforcedAttrs[key] = value
		}

		// enforced protocols
		for attr, protos := range tagdef.EnforcedProtocols {
			if newdef.EnforcedProtocols == nil {
				newdef.EnforcedProtocols = make(Protomap)
			}
			protoset := make(Protoset)
			for proto := range protos {
				protoset[proto] = struct{}{}
			}
			newdef.EnforcedProtocols[attr] = protoset
		}

		out[tag] = newdef
	}
	return out
}

func normalizeAttrKey(key string) string {
	return strings.ToLower(strings.Map(func(r rune) rune {
		if strings.IndexRune("\n\r\t />\"='\u0000", r) < 0 {
			return r
		}
		return -1
	}, key))
}

func normalizeProtocol(proto string) string {
	proto = strings.Map(func(r rune) rune {
		if strings.IndexRune("abcdefghijklmnopqrstuvwxyz0123456789+-.", r) >= 0 {
			return r
		}
		return -1
	}, strings.ToLower(proto))
	if proto != "" && []rune(proto)[0] >= 'a' && []rune(proto)[0] <= 'z' {
		return proto
	}
	return ""
}
