package main

import (
	//"github.com/pocketbase/pocketbase"
	//"github.com/pocketbase/pocketbase/apis"

	"net/url"

	"github.com/pocketbase/pocketbase/core"
	//"github.com/pocketbase/pocketbase/plugins/migratecmd"
	//. "maragu.dev/gomponents"
	//. "maragu.dev/gomponents/html"

	_ "kron/migrations"
	"kron/models"
	"kron/views"
)

func gDashboard(r *core.RequestEvent) error {
	if (r.Auth == nil) {
		r.Redirect(302, "/login")
	}
	html := views.DashboardHome()
	return html.Render(r.Response)
}

func gJobs(r *core.RequestEvent) error {
	if (r.Auth == nil) {
		r.Redirect(302, "/login")
	}
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
	if (r.Auth == nil) {
		r.Redirect(302, "/login")
	}
	e := r.Request.ParseForm()
	if e != nil {
		return r.Error(400, "Bad request", e)
	}
	requiredFields := [3]string{"name", "target", "method"}
	for _, field := range requiredFields {
		println(field + ": checking")
		if !r.Request.Form.Has(field) {
			return r.Error(400, "missing field "+field, nil)
		}
	}

	collection, err := r.App.FindCollectionByNameOrId("jobs")
	if err != nil {
		return err
	}

	var job models.Job

	job.Name = r.Request.Form.Get("name")

	_, err = url.Parse(r.Request.Form.Get("target"))
	if err != nil {
		return err
	}
	job.Target = r.Request.Form.Get("target")

	job.Request.Method = r.Request.Form.Get("method")
	job.Request.Body = r.Request.Form.Get("request")
	job.Expected_response = r.Request.Form.Get("expected_response")
	job.User_id = r.Auth.Id
	job.Schedule = r.Request.Form.Get("schedule")
	
	var record = core.NewRecord(collection)
	record.Set("user_id", job.User_id)
	record.Set("name", job.Name)
	record.Set("target", job.Target)
	record.Set("request", job.Request)
	record.Set("expected_response", job.Expected_response)
	record.Set("schedule", job.Schedule)

	if err := r.App.Save(record); err != nil {
		return err
	}
	
	err = reloadJobs(r.App)

	r.String(200, "Job added successfully")

	return nil

}
