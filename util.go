package gsoup

func cloneWhitelist(in whitelist) (out whitelist) {
	for tag, tagdef := range in {
		newdef := Tagdef{Tag: tagdef.Tag, AllowedAttrs: make(Attrset)}
		for attr := range tagdef.AllowedAttrs {
			newdef.AllowedAttrs[attr] = struct{}{}
		}
		out[tag] = newdef
	}
	return out
}
