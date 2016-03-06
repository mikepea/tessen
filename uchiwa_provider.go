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

func GetFilteredListOfUchiwaEvents(query Query, data *interface{}) []interface{} {
	results := make([]interface{}, 0)
	if data == nil {
		return results
	}
	templateName := "uchiwa_event_list"
	if query.Template != "" {
		templateName = query.Template
	}
	template := GetTemplate(templateName)
	if template == "" {
		template = GetTemplate("uchiwa_event_list")
	}
	b, err := json.Marshal(*data)
	if err != nil {
		log.Errorf("%s", err)
		return results
	}

	cmd := exec.Command("jq", fmt.Sprintf(".[] | select( %s ) | ._id", query.Filter))
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
			if event["_id"].(string) == id {
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

func FetchUchiwaEvent(id string, source *Source) interface{} {
	eventData := source.CachedData.([]interface{})
	for _, ev := range eventData {
		event := ev.(map[string]interface{})
		if event["_id"].(string) == id {
			return event
		}
	}
	return nil
}

func FetchUchiwaEvents(s *Source) ([]interface{}, error) {
	var contents []byte
	var err error
	log.Debugf("FetchUchiwaEvents: %s %s", s.Name, s.Endpoint)
	if s.Endpoint[0] == '/' {
		contents, err = getUchiwaResultsFromFile(s)
	} else {
		contents, err = getUchiwaResultsFromUchiwa(s)
	}

	if err != nil {
		return nil, err
	}

	data := make([]interface{}, 0)
	dec := json.NewDecoder(strings.NewReader(string(contents)))
	dec.Decode(&data)
	return data, nil

}

func getUchiwaResultsFromFile(s *Source) (contents []byte, err error) {
	contents, err = ioutil.ReadFile(s.Endpoint)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func getUchiwaResultsFromUchiwa(s *Source) (contents []byte, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/events", s.Endpoint))
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
