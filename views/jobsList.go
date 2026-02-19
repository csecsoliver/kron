package views

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func JobsList(names []string, urls []string) Node {
	items := []Node{}
	for i := range names {
		items = append(items, Li(
			Text(names[i]+" - "+urls[i]),
		))
	}
	return Ul(
		items...,
	)
}
