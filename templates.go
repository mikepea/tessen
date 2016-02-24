package tessen

var all_templates = map[string]string{
	"debug":      default_debug_template,
	"list":       default_list_template,
	"event_view": default_event_view_template,
	"help":       default_help_template,
}

//   issued:      [{{ dateFormat "2006-01-02T15:04" .check.issued }}](fg-blue)
const (
	default_debug_template      = "{{ . | toJson}}\n"
	default_list_template       = `[{{ .key | printf "%-12s"}}](fg-red)  [{{ if .fields.assignee }}{{ .fields.assignee.name | printf "%-10s" }}{{else}}{{"Unassigned"| printf "%-10s" }}{{end}} ](fg-blue) [{{ .fields.status.name | printf "%-12s"}}](fg-blue) [{{ dateFormat "2006-01-02" .fields.created }}](fg-blue)/[{{ dateFormat "2006-01-02T15:04" .fields.updated }}](fg-green)  {{ .fields.summary | printf "%-75s"}}`
	default_event_view_template = `
Event:
  _id:          [{{ ._id }}](fg-blue)
  acknowledged: [{{ .acknowledged }}](fg-blue)

Client:
  name:         [{{ .client.name }}](fg-blue)
  instance_id:  [{{ .client.instance_id }}](fg-blue)
{{ if .client.tags }}
  fqdn:         [{{ .client.tags.FQDN }}](fg-blue)
  ecosystem:    [{{ .client.tags.Ecosystem }}](fg-blue)
  region:       [{{ .client.tags.Region }}](fg-blue)
  display_name: [{{ index . "client" "tags" "Display Name" }}](fg-blue)
{{ end }}

Check:
  name:        [{{ .check.name }}](fg-green)
  command:     [{{ .check.command }}](fg-green)
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
)
