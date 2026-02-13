package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func PublicLayout(title string, body Node) Node {
	return El("HTML", Lang("en"), Head(
		Meta(Charset("UTF-8")),
		Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0")),
		TitleEl(Text(title)),
		Script(Src("/static/htmx.min.js")),
		Link(Rel("icon"), Href("/static/favicon.png")),
		Link(Rel("stylesheet"), Href("/static/styles.css")),
	),
		Body(
			Header(
				H1(Text("Kron")),
				Nav(Ul(
					Li(
						A(Href("/"), Text("Home")),
					),
					Li(
						A(Href("/login"), Text("Login/Register")),
					),
					Li(
						A(Href("https://github.com/csecsoliver/kron"), Target("blank"), Text("Github")),
					),
				)),
			),
			Main(body),
		),
	)
}
