package gsoup

import "golang.org/x/net/html/atom"

var simpleTextWhitelist = whitelist{
	atom.B:      T(atom.B),
	atom.Em:     T(atom.Em),
	atom.I:      T(atom.I),
	atom.Strong: T(atom.Strong),
	atom.U:      T(atom.U),
}

var basicWhitelist = whitelist{
	atom.A:          T(atom.A, "href").Enforce("rel", "nofollow"),
	atom.B:          T(atom.B),
	atom.Blockquote: T(atom.Blockquote, "cite"),
	atom.Br:         T(atom.Br),
	atom.Cite:       T(atom.Cite),
	atom.Dd:         T(atom.Dd),
	atom.Dl:         T(atom.Dl),
	atom.Dt:         T(atom.Dt),
	atom.Em:         T(atom.Em),
	atom.I:          T(atom.I),
	atom.Li:         T(atom.Li),
	atom.Ol:         T(atom.Ol),
	atom.P:          T(atom.P),
	atom.Pre:        T(atom.Pre),
	atom.Q:          T(atom.Q, "cite"),
	atom.Small:      T(atom.Small),
	atom.Span:       T(atom.Span),
	atom.Strike:     T(atom.Strike),
	atom.Strong:     T(atom.Strong),
	atom.Sub:        T(atom.Sub),
	atom.Sup:        T(atom.Sup),
	atom.U:          T(atom.U),
	atom.Ul:         T(atom.Ul),
}

var deleteChildrenSet = Tagset{
	atom.Applet:   struct{}{},
	atom.Area:     struct{}{},
	atom.Audio:    struct{}{},
	atom.Base:     struct{}{},
	atom.Basefont: struct{}{},
	atom.Br:       struct{}{},
	atom.Canvas:   struct{}{},
	atom.Col:      struct{}{},
	atom.Colgroup: struct{}{},
	atom.Datalist: struct{}{},
	atom.Embed:    struct{}{},
	atom.Frame:    struct{}{},
	atom.Frameset: struct{}{},
	atom.Head:     struct{}{},
	atom.Hr:       struct{}{},
	atom.Iframe:   struct{}{},
	atom.Img:      struct{}{},
	atom.Input:    struct{}{},
	atom.Keygen:   struct{}{},
	atom.Link:     struct{}{},
	atom.Map:      struct{}{},
	atom.Menu:     struct{}{},
	atom.Meta:     struct{}{},
	atom.Noframes: struct{}{},
	atom.Noscript: struct{}{},
	atom.Object:   struct{}{},
	atom.Param:    struct{}{},
	atom.Progress: struct{}{},
	atom.Rp:       struct{}{},
	atom.Script:   struct{}{},
	atom.Source:   struct{}{},
	atom.Style:    struct{}{},
	atom.Textarea: struct{}{},
	atom.Track:    struct{}{},
	atom.Video:    struct{}{},
	atom.Wbr:      struct{}{},
}

var preserveChildrenSet = Tagset{
	atom.Html: struct{}{},
	atom.Body: struct{}{},
}
