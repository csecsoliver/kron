package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func BaseLayout(title string, children ...Node) Node {
	return Doctype(
		El("HTML",
			Lang("en"),
			Head(
				Meta(Charset("UTF-8")),
				Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0")),
				TitleEl(Text(title)),
				Script(Src("/static/htmx.min.js")),
				Link(Rel("icon"), Href("/static/favicon.png")),
				Link(Rel("stylesheet"), Href("/static/normalize.css")),
				Link(Rel("stylesheet"), Href("/static/styles.css")),
			),
			Body(
				children...,
			),
		),
	)
}
