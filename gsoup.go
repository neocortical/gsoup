package gsoup

import "golang.org/x/net/html/atom"

// NewEmptyCleaner creates a cleaner with no allowed tags
func NewEmptyCleaner() Cleaner {
	return &cleaner{w: whitelist{}}
}

// NewBasicCleaner creates a Cleaner with the default basic whitelist
// This whitelist mirror's Jsoup's basic whitelist
func NewBasicCleaner() Cleaner {
	return &cleaner{w: cloneWhitelist(basicWhitelist)}
}

// NewBasicCleanerWithImages creates a Cleaner with the default basic whitelist
// that also allows <img> with http or https protocols.
// This whitelist mirror's Jsoup's basic whitelist with images.
func NewBasicCleanerWithImages() Cleaner {
	return &cleaner{w: cloneWhitelist(basicWhitelistWithImages)}
}

func NewRelaxedCleaner() Cleaner {
	return &cleaner{w: cloneWhitelist(relaxedWhitelist)}
}

// T is a shorthand method for creating a new Tagdef
func T(tag atom.Atom, attrs ...string) (def *Tagdef) {
	def = &Tagdef{Tag: tag}
	def.AllowedAttrs = make(Attrset)
	for _, attr := range attrs {
		def.AllowedAttrs[normalizeAttrKey(attr)] = struct{}{}
	}
	return def
}
