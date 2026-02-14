package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DashboardLayout(title string, children ...Node) Node {
	return Div(Class("dashboard"),
		Header(H1(Text(title))),
		Main(children...),
	)
}