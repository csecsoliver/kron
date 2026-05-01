package main

import (
	//"github.com/pocketbase/pocketbase"
	//"github.com/pocketbase/pocketbase/apis"

	"net/url"
	"regexp"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"

	//"github.com/pocketbase/pocketbase/plugins/migratecmd"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

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
	records, err := r.App.FindRecordsByFilter("jobs","user_id = {:userid}", "-created_at", 0,0, dbx.Params{"userid": r.Auth.Id})
	if err != nil {
		return err
	}
	jobs, err := jobRecordsToStructs(records)
	if err != nil {
		return err
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
		r.String(400, "Bad request")
		return e
	}
	requiredFields := [4]string{"name", "target", "method", "schedule"}
	for _, field := range requiredFields {
		println(field + ": checking")
		if !r.Request.Form.Has(field) {
			r.String(200, "missing field: "+field)
			return nil
		}
	}

	collection, err := r.App.FindCollectionByNameOrId("jobs")
	if err != nil {
		r.String(200, "Server not set up")
		return err
	}

	var job models.Job

	job.Name = r.Request.Form.Get("name")

	_, err = url.Parse(r.Request.Form.Get("target"))
	if err != nil {
		r.String(200, "Invalid url")
		return err
	}
	job.Target = r.Request.Form.Get("target")

	job.Request.Method = r.Request.Form.Get("method")
	job.Request.Body = r.Request.Form.Get("request")
	job.Expected_response = r.Request.Form.Get("expected_response")
	job.User_id = r.Auth.Id
	job.Schedule = r.Request.Form.Get("schedule")

	cron_regexp, err :=  regexp.Compile(`(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})`)
	if err != nil {
		r.String(200, "Server error")
		return err
	}
	if !cron_regexp.MatchString(job.Schedule) {
		r.String(200, "Invalid cron expression")
		return nil
	}
	
	var record = core.NewRecord(collection)
	record.Set("user_id", job.User_id)
	record.Set("name", job.Name)
	record.Set("target", job.Target)
	record.Set("request", job.Request)
	record.Set("expected_response", job.Expected_response)
	record.Set("schedule", job.Schedule)

	if err := r.App.Save(record); err != nil {
		r.String(200, err.Error())
		return err
	}

	err = reloadJobs(r.App)

	records, err := r.App.FindAllRecords("jobs")
	if err != nil {
		r.String(200, err.Error())
		return err
	}
	jobs, err := jobRecordsToStructs(records)
	if err != nil {
		r.String(200, err.Error())
		return err
	}
	return Span(
		Text("Job added successfully"),
		Div(
			Attr("hx-swap-oob", "innerHTML"),
			ID("jobslist"),
			views.JobsList(jobs),
		),
		).Render(r.Response)

	// r.String(200, "Job added successfully")
	// return nil

}

func gJob(r *core.RequestEvent) error {
	if (r.Auth == nil) {
		r.Redirect(302, "/login")
	}
	record, err := r.App.FindRecordById("jobs", r.Request.PathValue("id"))
	if err != nil {
		r.String(200, err.Error())
		return err
	}
	job, err := jobRecordToStruct(record)
	if err != nil {
		r.String(200, err.Error())
		return err
	}
	if job.User_id != r.Auth.Id {
		r.String(403, "Forbidden")
		return nil
	}
	records, err := r.App.FindRecordsByFilter("status_logs","job_id = {:jobid}", "-created_at", 0,0, dbx.Params{"jobid": job.Id})
	if err != nil {
		r.String(200, err.Error())
		return err
	}
	return views.JobDetails(job, records).Render(r.Response)
}
