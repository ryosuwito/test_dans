package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/parnurzeal/gorequest"
)

func GetJobs(description, location string, fullTime bool, page int) ([]Job, error) {
	request := gorequest.New()
	endpoint := "http://dev3.dansmultipro.co.id/api/recruitment/positions.json"

	params := url.Values{}
	if description != "" {
		params.Add("description", description)
	}
	if location != "" {
		params.Add("location", location)
	}
	if fullTime {
		params.Add("full_time", "true")
	}

	_, body, errs := request.Get(endpoint).Query(params).End()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	fmt.Println(body)
	var jobs []Job
	if err := json.Unmarshal([]byte(body), &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func GetJobDetail(id string) (*Job, error) {
	request := gorequest.New()
	endpoint := fmt.Sprintf("http://dev3.dansmultipro.co.id/api/recruitment/positions/%s", id)

	_, body, errs := request.Get(endpoint).End()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	var job Job
	if err := json.Unmarshal([]byte(body), &job); err != nil {
		return nil, err
	}

	return &job, nil
}
