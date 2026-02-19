package main

import (
	"github.com/pocketbase/pocketbase/core"
	"encoding/json"
	"kron/models"
)


func jobRecordToStruct(record *core.Record){
	var job models.Job
	job.name = record.GetString("name")
	job.target = record.GetString("target")
	json.Unmarshal([]byte(record.GetString("request")), &job.request)
	job.expected_response = record.GetString("expected_response")
	job.schedule = record.GetString("schedule")
	
}