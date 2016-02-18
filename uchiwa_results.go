package tessen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
