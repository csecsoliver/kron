package models

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"github.com/pocketbase/pocketbase/core"
	"github.com/agnivade/levenshtein"
)


type Request struct {
	Method  string            `json:"method"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}
type Job struct {
	Id                string	  `json:"id"`
	User_id           string  `json:"user_name"`
	Name              string  `json:"name"`
	Target            string  `json:"target"`
	Request           Request `json:"request"`
	Expected_response string  `json:"expected_response"`
	Schedule          string  `json:"schedule"`
}

type IJob interface {
	RegisterCron(app core.App)
}

func (j Job) RegisterCron(app core.App) error {
	println("refreshing cron with job id  " + j.Id+" and name "+j.Name)
	
	err := app.Cron().Add(j.Id, j.Schedule, func() {
		body := bytes.NewBufferString(j.Request.Body)
		collection, err := app.FindCollectionByNameOrId("status_logs")
		if err!=nil {
			panic(err)
		}
		request, err := http.NewRequest(j.Request.Method, j.Target, body)
		record := core.NewRecord(collection)
		record.Set("job_id", j.Id)
		
		if err != nil {
			record.Set("successful", false)
			record.Set("error", err)
			err := app.Save(record)
			if err != nil {
				panic(err)
			}
			return 
		} 
		for i, k := range j.Request.Headers {
			request.Header.Add(i, k)
		}
		client := http.Client{}
		response, err := client.Do(request)
		if err != nil {
			record.Set("successful", false)
			record.Set("error", err)
			err := app.Save(record)
			if err != nil {
				panic(err)
			}
			return 
		}
		defer response.Body.Close()
		record.Set("status_code", response.StatusCode)
		record.Set("headers", response.Header)
		
		// Partial Source - https://stackoverflow.com/a/9649061
		// Posted by Stephen Weinberg, modified by community. See post 'Timeline' for change history
		// Retrieved 2026-04-30, License - CC BY-SA 4.0
		buf := new(strings.Builder)
		_, err = io.Copy(buf, response.Body)
		if err != nil {
			record.Set("successful", false)
			record.Set("error", err)
			err := app.Save(record)
			if err != nil {
				panic(err)
			}
			
			return 
		}
		
		record.Set("body", buf.String())
		if (strCmp(j.Expected_response, buf.String(), 0) || len(j.Expected_response) == 0){
			record.Set("successful", true)
		} else {
			
			record.Set("successful", false)
		}
		
		println(record.GetBool("successful"))
		err = app.Save(record)
		if err != nil {
			record.Set("successful", false)
			record.Set("error", err)
			return 
		}
		
	})
	for _, job := range app.Cron().Jobs() {
		println("cron job found: "+job.Id())
	}
	return err
}
func strNormalize(s string) string {
	return strings.Join(strings.Fields(s), "")
}
func strCmp(one, two string, fuzz int) bool {
	if (levenshtein.ComputeDistance(strNormalize(one), strNormalize(two)) <= fuzz) {
		return true
	}
	return false
}
