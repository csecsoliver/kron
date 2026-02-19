package main

import (
	//"github.com/pocketbase/pocketbase"
	//"github.com/pocketbase/pocketbase/apis"

	"github.com/pocketbase/pocketbase/core"
	//"github.com/pocketbase/pocketbase/plugins/migratecmd"
	//. "maragu.dev/gomponents"
	//. "maragu.dev/gomponents/html"

	_ "kron/migrations"
	"kron/views"
)

func gDashboard(r *core.RequestEvent) error {
	html := views.DashboardHome()
	return html.Render(r.Response)
}

func gJobs(r *core.RequestEvent) error {
	records, err := r.App.FindAllRecords("jobs")
	if err != nil {
		return err
	}
	names, urls := []string{}, []string{}
	for _, record := range records {
		names = append(names, record.GetString("name"))
		urls = append(urls, record.GetString("url"))
	}
	html := views.JobsList(names, urls)
	return html.Render(r.Response)
}
