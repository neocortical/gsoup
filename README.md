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
cleaner := gsoup.NewBasicCleaner().AddTags(
		T(atom.Div, "id", "class"),
		T(atom.Canvas),
	)
cleaner = cleaner.RemoveTags(atom.P)

var markup = `<p>deleted</p><div id="foo" class="bar">also saved</div><canvas></canvas>`

doc, err := gsoup.NewBasicCleaner().Clean(markup)
// doc is a html.Node that will render '<div id="foo" class="bar">also saved</div><canvas></canvas>'

// new Cleaner with no allowed tags
cleaner = gsoup.NewEmptyCleaner()
```

## TODO

* Enforced attributes (e.g. rel="nofollow")
* Protocol enforcement for src, href attributes
* HTML transformers like in rgrove/sanitize
* ???

## Caveats

This package is in Alpha and may change. Comments, feature requests, bug reports, pull requests all welcome!

Version 0.2.0

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
