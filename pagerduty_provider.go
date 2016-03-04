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

func GetFilteredListOfPagerDutyEvents(query Query, data *interface{}) []interface{} {
	log.Debugf("GetFilteredListOfPagerDutyEvents: data: %q", *data)
	results := make([]interface{}, 0)
	if data == nil {
		return results
	}
	templateName := "pagerduty_incident_list"
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

	cmd := exec.Command("jq", fmt.Sprintf(".[] | select( %s ) | .id", query.Filter))
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
			if event["id"].(string) == id {
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

func FetchPagerDutyEvent(id string, source Source) interface{} {
	eventData := source.CachedData.([]interface{})
	log.Debugf("FetchPagerDutyEvent: eventData: %q", eventData)
	for _, ev := range eventData {
		log.Debugf("FetchPagerDutyEvent: ev: %q", ev)
		event := ev.(map[string]interface{})
		if event["id"].(string) == id {
			return event
		}
	}
	return nil
}

func FetchPagerDutyEvents(endpoint string) ([]interface{}, error) {
	var contents []byte
	var err error
	if endpoint[:7] == "file://" {
		contents, err = getPagerDutyResultsFromFile(endpoint[7:])
	} else {
		contents, err = getPagerDutyResultsFromPagerDuty(endpoint)
	}
	log.Debugf("FetchPagerDutyEvents: contents: %q", contents)

	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	dec := json.NewDecoder(strings.NewReader(string(contents)))
	dec.Decode(&data)
	incidents := data["incidents"].([]interface{})
	log.Debugf("FetchPagerDutyEvents: incidents: %q", incidents)
	return incidents, nil

}

func getPagerDutyResultsFromFile(file string) (contents []byte, err error) {
	contents, err = ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func getPagerDutyResultsFromPagerDuty(endpoint string) (contents []byte, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/incidents", endpoint))
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
