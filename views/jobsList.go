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
			Text(jobs[i].Name + " - " + jobs[i].Target),
		))
	}
	return Ul(
		items...,
	)
}
