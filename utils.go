package main

import (
	"errors"
	"fmt"
	"kron/models"

	"github.com/pocketbase/pocketbase/core"
)

func jobRecordToStruct(record *core.Record) (models.Job, error) {
	var job models.Job
	job.Id = record.GetString("id")
	job.User_id = record.GetString("user_id")
	job.Name = record.GetString("name")
	job.Target = record.GetString("target")
	job.Expected_response = record.GetString("expected_response")
	job.Schedule = record.GetString("schedule")
	println("job id: " + job.Id)
	println("name: " + job.Name)
	println("schedule: " + job.Schedule)
	err := record.UnmarshalJSONField("request", &job.Request)
	// err := json.Unmarshal([]byte(record.GetString("request")), &job.Request)
	if err != nil {
		return job, err
	}
	return job, nil
}
func jobRecordsToStructs(records []*core.Record) ([]models.Job, error) {
	var jobs []models.Job
	for _, record := range records {
		job, err := jobRecordToStruct(record)
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
func validate_method(method string) error {
	switch method {
	case "GET", "POST", "PUT", "DELETE":
		return nil
	default:
		return errors.New("invalid method")
	}
}
func register_delete_old_by_job(app core.App) error {
	return app.Cron().Add("delete_old_by_job", "0 * * * *", func() {
		q := app.DB().NewQuery(`WITH ranked AS (
			SELECT id,
				ROW_NUMBER() OVER (PARTITION BY job_id ORDER BY created_at DESC) AS rn
			FROM status_logs
		)
		DELETE FROM status_logs
		WHERE id IN (SELECT id FROM ranked WHERE rn > 60)
		`)
		r, err := q.Execute()
		if err != nil {
			println("Error deleting old status logs: " + err.Error())
		} else {
			rowsAffected, err := r.RowsAffected()
			if err != nil {
				println("Error getting rows affected: " + err.Error())
			}
			println("Deleted old status logs: " + fmt.Sprintf("%d",rowsAffected))
		}
	})
}
