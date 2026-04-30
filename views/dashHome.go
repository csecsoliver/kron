package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DashboardHome() Node {
	return DashboardLayout("Dashboard", Div(
		ID("dHome"),
		Div(
			ID("newjob"),
			Form(
				Attr("hx-post", "/dash/jobs"),
				Attr("hx-target", "#newjobres"),
				Attr("hx-swap", "innerHTML"),
				
				Label(For("name"), Text("Name")),
				Input(Type("text"), Name("name"), Placeholder("Job name"), Required()),
				
				Label(For("target"), Text("Target URL")),
				Input(Type("text"), Name("target"), Placeholder("URL to check"), Required()),
				
				Label(For("method"), Text("Method")),
				Select(Name("method"), Required(),
					Option(Value("GET"), Text("GET")),
					Option(Value("POST"), Text("POST")),
					Option(Value("PUT"), Text("PUT")),
					Option(Value("DELETE"), Text("DELETE")),
				),
				
				Label(For("request"), Text("Request body (optional)")),
				Textarea(Name("request"), Placeholder("Request body")),
				
				Label(For("expected_response"), Text("Expected Response (optional)")),
				Textarea(Name("expected_response"), Placeholder("Expected response body")),
				
				Label(For("schedule"), Text("Schedule (cron format)")),
				Input(Type("text"), ID("schedule"), Name("schedule")),
				
				Button(Type("submit"), Text("Create new job")),
			),
			Div(
				ID("newjobres"),
			),
		),
		Div(
			Button(
				Attr("hx-get", "/dash/jobs"),
				Attr("hx-target", "#jobslist"),
				Text("Reload jobs"),
			),
			Div(
				ID("jobslist"),
				Attr("hx-get", "/dash/jobs"),
				Attr("hx-target", "this"),
				Attr("hx-swap", "innerHTML"),
				Attr("hx-trigger", "load"),
			),
		),
		Div(
			Div(
				ID("jobdetails"),
				
			),
		),
	))
}
