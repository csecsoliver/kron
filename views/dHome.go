package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)
func DashboardHome() Node {
	return DashboardLayout("Dashboard", Div(
		Div(
			ID("newjob"),
			Form(
				Attr("hx-post", "/dash/jobs/new"),
				Attr("hx-target", "#newjobres"),
				Attr("hx-swap", "innerHTML"),
				Input(Type("text"), Name("name"), Placeholder("Job name")),
			),
			Div(
				ID("newjobres"),	
			),
		),
		Div(
			Div(
				ID("jobslist"),
				Attr("hx-get", "/dash/jobs/list"),
				Attr("hx-target", "this"),
				Attr("hx-swap", "innerHTML"),
				Attr("hx-trigger", "reload-jobs from:body"),
			),
			Button(
				Attr("hx-post", "/dash/jobs/list"),
				Attr("hx-target", "#jobslist"),
				Text("Reload jobs"),
			),
		),
	))
}
