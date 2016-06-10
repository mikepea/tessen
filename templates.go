package tessen

var all_templates = map[string]string{
	"debug":                   default_debug_template,
	"uchiwa_event_list":       uchiwa_event_list_template,
	"uchiwa_team_event_list":  uchiwa_team_event_list_template,
	"uchiwa_event_view":       uchiwa_event_view_template,
	"pagerduty_incident_list": pagerduty_incident_list_template,
	"pagerduty_incident_view": default_debug_template,
	"chronos_job_list":        chronos_job_list_template,
	"chronos_job_view":        default_debug_template,
	"help":                    default_help_template,
	"event_http_dashboard_home": default_event_http_dashboard_home_template,
}

//   issued:      [{{ dateFormat "2006-01-02T15:04" .check.issued }}](fg-blue)
const (
	default_debug_template           = "{{ . | toJson}}\n"
	uchiwa_event_list_template       = `{{ colorizedSensuStatus .check.status | printf "%-6s"}}  [{{ .check.name | printf "%-40s" }}](fg-green)  {{ .client.name }}`
	uchiwa_team_event_list_template  = `{{ colorizedSensuStatus .check.status | printf "%-6s"}}  [{{ .check.team | printf "%-20s" }}](fg-blue)  [{{ .check.name | printf "%-40s" }}](fg-green)  {{ .client.name }}`
	pagerduty_incident_list_template = `{{ .id }}  {{.status | printf "%-12s"}}  [{{ .escalation_policy.name | printf "%-40s" }}](fg-green)  [{{if .assigned_to_user}}{{ .assigned_to_user.email | printf "%-20s" }}{{else}}{{ "UNASSIGNED" | printf "%-40s" }}{{end}}](fg-green)  {{ .trigger_summary_data.description }}`
	chronos_job_list_template        = `[{{ .name | printf "%-60s" }}](fg-green) | {{ (print .cpus .mem .disk) | printf "%-15s" }} | {{.schedule | printf "%-30s" }} {{.epsilon}}`
	uchiwa_event_view_template       = `
Event:
  _id:          [{{ ._id }}](fg-blue)
  acknowledged: [{{ .acknowledged }}](fg-blue)
  occurrences:  [{{ .occurrences }}](fg-blue)

Client:
  name:         [{{ .client.name }}](fg-green)
  instance_id:  [{{ .client.instance_id }}](fg-blue)
{{ if .client.tags }}
  fqdn:         [{{ .client.tags.FQDN }}](fg-blue)
  ecosystem:    [{{ .client.tags.Ecosystem }}](fg-green)
  region:       [{{ .client.tags.Region }}](fg-blue)
  display_name: [{{ index . "client" "tags" "Display Name" }}](fg-blue)
{{ end }}

Check:
  name:        [{{ .check.name }}](fg-green)
  command:     [{{ .check.command }}](fg-green)
  interval:    [{{ .check.interval }}](fg-blue)
  issued:      [{{ .check.issued }}](fg-blue)
  team:        [{{ .check.team }}](fg-blue)
  project:     [{{ .check.project }}](fg-blue)
  status:      {{ colorizedSensuStatus .check.status }}

  runbook:     [{{ .check.runbook }}](fg-blue)

  page:        [{{ .check.page }}](fg-blue)
  ticket:      [{{ .check.ticket }}](fg-blue)
{{ if .check.history }}
  history:    {{range .check.history}} {{ colorizedSensuStatus . }}{{end}}
{{ end }}

Output:

  {{ indent 2 .check.output }}

`
	default_help_template = `
[Quick reference for tessen](fg-white)

[Actions:](fg-blue)

    <enter>      - select item
    h            - show help page

[Commands (a'la vim/tig):](fg-blue)

    :query {JQ boolean expression} - display filtered results
    :help                          - show help page
    :<up>                          - select previous command
    :quit or :q                    - quit

[Navigation:](fg-blue)

    up/k         - previous line
    down/j       - next line
    C-f/<space>  - next page
    C-b          - previous page
    }            - next paragraph/section/fast-move
    {            - previous paragraph/section/fast-move
    n            - next search match
    g            - go to top of page
    G            - go to bottom of page
    q            - go back / quit
    C-c/Q        - quit

[Notes:](fg-blue)

    Learning JQ is highly recommended, particularly boolean expressions:

        https://stedolan.github.io/jq/manual/

`
	default_event_http_dashboard_home_template = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta http-equiv="refresh" content="10">
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
		<title>Tessen Dashboard</title>

    <!-- Bootstrap -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">


  <style>
  body {
    background-color: #002b36;
    color: rgba(255, 255, 255, 1);
  }
  .query-overview {
    font-size: xx-large;
    text-align: center;
  }
  .query-overview .critical {
    border-radius: 25px;
    background-color: #dc322f;
    color: rgba(255, 255, 255, 1);
  }
  .query-overview .warning {
    border-radius: 25px;
    background-color: #b58900;
    color: rgba(255, 255, 255, 1);
  }
  .query-overview .unknown {
    border-radius: 25px;
    background-color: #268bd2;
    color: rgba(255, 255, 255, 1);
  }
  .query-overview .ok {
    border-radius: 25px;
    background-color: #859900;
    color: rgba(255, 255, 255, 1);
  }
  .query-overview .count .critical {
    border-radius: 0px;
    display: inline-block;
    font-size: xx-large;
    padding: 3px;
  }
  .query-overview .count .warning {
    border-radius: 0px;
    display: inline-block;
    font-size: xx-large;
    padding: 3px;
  }
  .query-overview .count .unknown {
    border-radius: 0px;
    display: inline-block;
    padding: 3px;
  }
  </style>

  </head>
  <body>

		<div class="container">
      <h1>My Dashboard</h1>
    </div>

		<div class="container">
      <div class="row">
	    {{range .}}
			  <div class="col-md-4">
					<div class="query-overview">
            {{if gt .CritCount 0 }}
						<div class="critical">
            {{else if gt .WarnCount 0 }}
						<div class="warning">
            {{else if gt .UnknCount 0 }}
						<div class="unknown">
            {{else}}
						<div class="ok">
            {{end}}
							<h2>{{ .Name }}</h2>
              <div class="count">
              <div class="critical">{{ .CritCount }}</div>
              <div class="warning">{{ .WarnCount }}</div>
              <div class="unknown">{{ .UnknCount }}</div>
            </div>
						</div>
					</div>
			  </div>
	      {{end}}
      </div>
		</div>

  </body>
</html>
`
)
