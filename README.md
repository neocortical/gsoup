# Gsoup: An HTML sanitizer for Go

## Description

Gsoup enables the sanitization of untrusted HTML markup for use in blogs, comments, etc. It is similar to Java's Jsoup or Ruby's Sanitizer (https://github.com/rgrove/sanitize). Go's x/net/html package is used under the hood for HTML parsing.

[![Build Status](https://travis-ci.org/neocortical/gsoup.svg?branch=master)](https://travis-ci.org/neocortical/gsoup) [![Coverage](http://gocover.io/_badge/github.com/neocortical/gsoup)](http://gocover.io/github.com/neocortical/gsoup) [![GoDoc](https://godoc.org/github.com/neocortical/gsoup?status.svg)](https://godoc.org/github.com/neocortical/gsoup)

## Installation

`go get github.com/neocortical/gsoup`

`import "github.com/neocortical/gsoup"`

## Basic Use

```go
var markup = `<p onclick="alert('XSSed!')">save me</p><div>delete me?</div>`

doc, err := gsoup.NewBasicCleaner().Clean(markup)
// doc is a html.Node that will render '<p>save me</p>'

cleaned, err := gsoup.NewBasicCleaner().PreserveChildren().Clean(markup)
// cleaned is a html.Node that will render '<p>save me</p>delete me?'

```
## Custom Use

```go
// new Cleaner with no allowed tags
cleaner = gsoup.NewEmptyCleaner()

// completely customize allowed tags
cleaner := gsoup.NewEmptyCleaner().AddTags(
		T(atom.Div, "id", "class"),
		T(atom.Canvas),
		T(atom.P),
	)

// RemoveTags is a factory method just like AddTags
cleaner = gsoup.NewBasicCleaner().RemoveTags(atom.P)

// EnforceProtocols enforces both the specified protocols and also valid URLs
// attributes with values that do not meet these requirements will be removed
cleaner = gsoup.NewEmptyCleaner().AddTags(
		T(atom.A, "href").EnforceProtocols("href", "http", "https", "mailto"),
	)

```

## TODO

* HTML transformers like in rgrove/sanitize
* CSS value sanitation?
* Even more tests for malicious vectors

## Caveats

This package is in Alpha and may change. Comments, feature requests, bug reports, pull requests all welcome!

Version 0.5.0

## License

The MIT License (MIT)

Copyright (c) 2015 Nathan Smith

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
