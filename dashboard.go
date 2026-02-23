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
	"kron/models"
	
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
	jobs := []models.Job{}
	for _, record := range records {
		err, job := jobRecordToStruct(record)
		if err != nil {
			return err
		}
		jobs = append(jobs, job)
	}
	html := views.JobsList(jobs)
	return html.Render(r.Response)
}
