package main

import (
	"kron/models"
	"github.com/pocketbase/pocketbase/core"
	"encoding/json"
)
func jobRecordToStruct(record *core.Record) (error, models.Job )	{
	var job models.Job
	job.Name = record.GetString("name")
	job.Target = record.GetString("target")
	job.Expected_response = record.GetString("expected_response")
	job.Schedule = record.GetString("schedule")
	err := json.Unmarshal([]byte(record.GetString("request")), &job.Request)
	if err != nil {
		return err, job
	}
	return nil, job
}
