// Package gsoup provides HTML sanitization functionality on top of Go's html package
//
// Parsing of input is handled by the html package. gsoup's Clean
// method cleans the resulting DOM, based on a whitelist of tags and the attributes
// and attribute protocols they are allowed to contain. Users can choose from a preconfigured whitelist,
// modify an existing whitelist, or generate their own.
//
// Unlike Java's Jsoup library, cleaning is done in-place and the structured DOM
// is returned, allowing for subsequent structured manipuation.
//
// This package is currently in alpha
package gsoup
