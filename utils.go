package main

import (
	"encoding/json"
	"kron/models"

	"github.com/pocketbase/pocketbase/core"
)

func jobRecordToStruct(record *core.Record) (models.Job, error) {
	var job models.Job
	job.Name = record.GetString("name")
	job.Target = record.GetString("target")
	job.Expected_response = record.GetString("expected_response")
	job.Schedule = record.GetString("schedule")
	err := json.Unmarshal([]byte(record.GetString("request")), &job.Request)
	if err != nil {
		return job, err
	}
	return job, nil
}
