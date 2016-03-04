package tessen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func FetchUchiwaEvent(id string, source Source) interface{} {
	eventData := source.CachedData.([]interface{})
	for _, ev := range eventData {
		event := ev.(map[string]interface{})
		if event["_id"].(string) == id {
			return event
		}
	}
	return nil
}

func FetchUchiwaEvents(endpoint string) ([]interface{}, error) {
	var contents []byte
	var err error
	log.Debugf("FetchUchiwaEvents: %s", endpoint[7:])
	if endpoint[:7] == "file://" {
		contents, err = getUchiwaResultsFromFile(endpoint[7:])
	} else {
		contents, err = getUchiwaResultsFromUchiwa(endpoint)
	}

	if err != nil {
		return nil, err
	}

	data := make([]interface{}, 0)
	dec := json.NewDecoder(strings.NewReader(string(contents)))
	dec.Decode(&data)
	return data, nil

}

func getUchiwaResultsFromFile(file string) (contents []byte, err error) {
	contents, err = ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func getUchiwaResultsFromUchiwa(endpoint string) (contents []byte, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/events", endpoint))
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
