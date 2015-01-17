package gsoup

import "golang.org/x/net/html/atom"

// NewEmptyCleaner creates a cleaner with no allowed tags
func NewEmptyCleaner() Cleaner {
	return &cleaner{w: whitelist{}}
}

// NewSimpleCleaner creates a Cleaner with the default simple whitelist
// This whitelist mirrors Jsoup's simple whitelist
func NewSimpleCleaner() Cleaner {
	return &cleaner{w: cloneWhitelist(simpleTextWhitelist)}
}

// NewBasicCleaner creates a Cleaner with the default basic whitelist
// This whitelist mirrors Jsoup's basic whitelist
func NewBasicCleaner() Cleaner {
	return &cleaner{w: cloneWhitelist(basicWhitelist)}
}

// NewBasicCleanerWithImages creates a Cleaner with the default basic whitelist
// that also allows <img> with http or https protocols.
// This whitelist mirrors Jsoup's basic whitelist with images.
func NewBasicCleanerWithImages() Cleaner {
	return &cleaner{w: cloneWhitelist(basicWhitelistWithImages)}
}

// NewRelaxedCleaner creates a Cleaner with the default relaxed whitelist
// This whitelist mirrors Jsoup's relaxed whitelist
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
