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

func FetchUchiwaEvents(endpoint string) ([]map[string]interface{}, error) {

	resp, err := http.Get(fmt.Sprintf("%s/events", endpoint))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := make([]map[string]interface{}, 0)
	dec := json.NewDecoder(strings.NewReader(string(contents)))
	dec.Decode(&data)
	return data, nil

}

func GetFilteredListOfEvents(filter string, eventData *[]map[string]interface{}) []string {
	results := make([]string, 0)
	b, err := json.Marshal(*eventData)
	if err != nil {
		log.Errorf("%s", err)
		return results
	}

	cmd := exec.Command("jq", fmt.Sprintf(".[] | select( %s ) | ._id", filter))
	cmd.Stdin = bytes.NewReader(b)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Errorf("%s", err)
		return results
	}

	results = strings.Split(out.String(), "\n")
	return results

}
