package tessen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func GetFilteredListOfChronosJobs(query Query, data *interface{}) []interface{} {
	log.Debugf("GetFilteredListOfChronosJobs: data: %q", *data)
	results := make([]interface{}, 0)
	if data == nil {
		return results
	}
	templateName := "chronos_job_list"
	if query.Template != "" {
		templateName = query.Template
	}
	template := GetTemplate(templateName)
	if template == "" {
		template = GetTemplate("event_list")
	}
	b, err := json.Marshal(*data)
	if err != nil {
		log.Errorf("%s", err)
		return results
	}

	cmd := exec.Command("jq", fmt.Sprintf(".[] | select( %s ) | .name", query.Filter))
	cmd.Stdin = bytes.NewReader(b)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Errorf("%s", err)
		return results
	}

	for _, v := range strings.Split(out.String(), "\n") {
		id := strings.Trim(v, "\"")
		for _, e := range (*data).([]interface{}) {
			event := e.(map[string]interface{})
			if event["name"].(string) == id {
				id := strings.Trim(v, "\"")
				buf := new(bytes.Buffer)
				RunTemplate(template, event, buf)
				results = append(results, QueryResult{id, buf.String(), event})
				continue
			}
		}
	}

	if len(results) == 0 {
		results = append(results, QueryResult{"", "No results found", nil})
	}
	return results
}

func FetchChronosJob(id string, source *Source) interface{} {
	eventData := source.CachedData.([]interface{})
	log.Debugf("FetchChronosJob: eventData: %q", eventData)
	for _, ev := range eventData {
		log.Debugf("FetchChronosJob: ev: %q", ev)
		event := ev.(map[string]interface{})
		if event["name"].(string) == id {
			return event
		}
	}
	return nil
}

func FetchChronosJobs(s *Source) ([]interface{}, error) {
	var contents []byte
	var err error
	if s.Endpoint[0] == '/' {
		contents, err = getChronosResultsFromFile(s)
	} else {
		contents, err = getChronosResultsFromAPI(s)
	}
	log.Debugf("FetchChronosJobs: contents: %q", contents)

	if err != nil {
		return nil, err
	}

	data := make([]interface{}, 0)
	dec := json.NewDecoder(strings.NewReader(string(contents)))
	dec.Decode(&data)
	log.Debugf("FetchChronosJobs: data: %q", data)
	return data, nil

}

func getChronosResultsFromFile(s *Source) (contents []byte, err error) {
	contents, err = ioutil.ReadFile(s.Endpoint)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func getChronosResultsFromAPI(s *Source) (contents []byte, err error) {
	queryParams := "?status=triggered,acknowledged"
	client := &http.Client{}
	url := fmt.Sprintf("%s/scheduler/jobs%s", s.Endpoint, queryParams)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil

}
