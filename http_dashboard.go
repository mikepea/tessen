package tessen

import (
	"bytes"
	"fmt"
	"net/http"
)

type DisplayWidget struct {
	Name      string
	CritCount int
	WarnCount int
	UnknCount int
}

func getDashboardQueries() (queries []Query) {
	opts := getOpts()
	if q := opts["queries"]; q != nil {
		qList := q.([]interface{})
		for _, v := range qList {
			q1 := v.(map[interface{}]interface{})
			q2 := make(map[string]string)
			for k, v := range q1 {
				switch k := k.(type) {
				case string:
					switch v := v.(type) {
					case string:
						q2[k] = v
					}
				}
			}
			if q2["dashboard"] == "true" {
				queries = append(queries, Query{q2["name"], q2["filter"], q2["template"], FindSourceByName(q2["source"])})
			}
		}
	}
	return queries
}

func countEvents(status float64, qr []interface{}) (count int) {
	for _, ev := range qr {
		data := ev.(QueryResult).Data
		if data == nil {
			continue
		}
		check := data.(map[string]interface{})["check"].(map[string]interface{})
		if check["status"].(float64) == status {
			count = count + 1
		}
	}
	return count
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	queries := getDashboardQueries()
	template := GetTemplate("event_http_dashboard_home")
	displayWidgets := make([]DisplayWidget, 0)
	for _, q := range queries {
		r := GetQueryResults(q)
		numCrit := countEvents(2, r)
		numWarn := countEvents(1, r)
		numUnkn := countEvents(3, r)
		displayWidgets = append(displayWidgets, DisplayWidget{q.Name, numCrit, numWarn, numUnkn})
	}
	buf := new(bytes.Buffer)
	RunTemplate(template, displayWidgets, buf)
	fmt.Fprintf(w, buf.String())
}

func StartHttpDashboard(listen string) error {
	http.HandleFunc("/", dashboardHandler)
	return http.ListenAndServe(listen, nil)
}
