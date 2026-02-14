package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func PublicLayout(title string, body Node) Node {
	return BaseLayout(
		title,
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
	)
}
