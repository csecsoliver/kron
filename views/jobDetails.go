package views

import (
	"kron/models"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/core"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func JobDetails(job models.Job, runs []*core.Record) Node {
	rows := []Node{
		Tr(
			Td(
				Text("Name"),
			),
			Td(
				Text(job.Name),
			),
		),
		Tr(
			Td(
				Text("Target"),
			),
			Td(
				Text(job.Target),
			),
		),
		Tr(
			Td(
				Text("Method"),
			),
			Td(
				Text(job.Request.Method),
			),
		),
		Tr(
			Td(
				Text("Schedule"),
			),
			Td(
				Text(job.Schedule),
			),
		),
		Tr(
			Td(
				Text("Request Body"),
			),
			Td(
				Text(job.Request.Body),
			),
		),
		Tr(
			Td(
				Text("Expected Response"),
			),
			Td(
				Text(job.Expected_response),
			),
		),
	}

	if len(runs) != 0 {
		lastRun := runs[0]
		loc, _ := time.LoadLocation("Europe/Budapest")
		t := lastRun.GetDateTime("created_at").Time()
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
		

		rows = append(rows, Tr(
			Td(
				Text("Last Run"),
			),
			Td(
				Text(t.Format(time.Layout)),
			),
		))
		rows = append(rows, Tr(
			Td(
				Text("Last Run Success"),
			),
			Td(
				Text(strconv.FormatBool(lastRun.GetBool("successful"))),
			),
		))
		rows = append(rows, Tr(
			Td(
				Text("Last Run Status Code"),
			),
			Td(
				Text(lastRun.GetString("status_code")),
			),
		))
		rows = append(rows, Tr(
			Td(
				Text("Last Run Response"),
			),
			Td(
				Text(lastRun.GetString("body")),
			),
		))
	}

	return Table(
		THead(),
		TBody(rows...),
	)
}
