package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func LoginPage(error string) Node {
	return PublicLayout(
		"Home",
		Div(
			H2(Text("Login")),
			H3(Class("error"),Text(error)),
			Form(
				Method("POST"),
				Div(
					Label(Text("Email: "), For("email")),
					Input(Type("email"), Name("email"), ID("email"), Required()),
				),
				Div(
					Label(Text("Password: "), For("password")),
					Input(Type("password"), Name("password"), ID("password"), Required()),
				),
				Button(Type("submit"), Text("Login/Register")),
			),
		),
	)
}
