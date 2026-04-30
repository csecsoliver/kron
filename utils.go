package main

import (
	"errors"
	"kron/models"

	"github.com/pocketbase/pocketbase/core"
)

func jobRecordToStruct(record *core.Record) (models.Job, error) {
	var job models.Job
	job.Id = record.GetString("id")
	job.Name = record.GetString("name")
	job.Target = record.GetString("target")
	job.Expected_response = record.GetString("expected_response")
	job.Schedule = record.GetString("schedule")
	println("job id: "+job.Id)
	println("name: "+job.Name)
	println("schedule: "+job.Schedule)
	err:=record.UnmarshalJSONField("request", &job.Request)
	// err := json.Unmarshal([]byte(record.GetString("request")), &job.Request)
	if err != nil {
		return job, err
	}
	return job, nil
}

func validate_method(method string) error {
	switch method {
	case "GET", "POST", "PUT", "DELETE":
		return nil
	default:
		return errors.New("invalid method")
	}
}