package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"kron/models"
)

func JobsList(jobs []models.Job) Node {
	items := []Node{}
	for i := range jobs {
		items = append(items, Li(
			Text(jobs[i].Name + " - " + jobs[i].Target + " "),
			Button(
				Attr("hx-get", "/dash/jobs/" + jobs[i].Id),
				Attr("hx-target", "#jobdetails"),
				Attr("hx-swap", "innerHTML"),
				Text("Details"),
			),
		))
	}
	return Ul(
		items...,
	)
}
