// Package gsoup provides HTML sanitization functionality on top of Go's html package
//
// Parsing of string-based input is handled by the html package. gsoup's Clean
// method cleans the resulting DOM, based on a whitelist of tags and the attributes
// they are allowed to contain. Users can choose from a preconfigured whitelist,
// modify an existing whitelist, or invent their own.
//
// Unlike Java's Jsoup library, cleaning is does in-place and the structured DOM
// is returned, allowing for subsequent structured manipuation.
//
// This package is currently in alpha
package gsoup

// Package version info
const VERSION = "0.1.0"
const MAJOR_VERSION = 0
const MINOR_VERSION = 1
const PATCH_VERSION = 0
