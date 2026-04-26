package main

import (
	//"github.com/pocketbase/pocketbase"
	//"github.com/pocketbase/pocketbase/apis"

	"github.com/pocketbase/pocketbase/core"
	//"github.com/pocketbase/pocketbase/plugins/migratecmd"
	//. "maragu.dev/gomponents"
	//. "maragu.dev/gomponents/html"

	_ "kron/migrations"
	"kron/models"
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
	var jobs []models.Job
	for _, record := range records {
		job, err := jobRecordToStruct(record)
		if err != nil {
			return err
		}
		jobs = append(jobs, job)
	}
	html := views.JobsList(jobs)
	return html.Render(r.Response)
}
func pJob(r *core.RequestEvent) error {
	e := r.Request.ParseForm()
	if e != nil {
		return r.Error(400, "Bad request", e)
	}
	requiredFields := [5]string{"name", "target", "method", "body", "expected"}
	for _, field := range requiredFields {
		if !r.Request.Form.Has(field) {
			return r.Error(400, "missing field"+field, nil)
		}
	}

	return nil

}
