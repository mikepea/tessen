package tessen

var all_templates = map[string]string{
	"debug":      default_debug_template,
	"list":       default_list_template,
	"event_view": default_event_view_template,
	"help":       default_help_template,
}

const (
	default_debug_template      = "{{ . | toJson}}\n"
	default_list_template       = `[{{ .key | printf "%-12s"}}](fg-red)  [{{ if .fields.assignee }}{{ .fields.assignee.name | printf "%-10s" }}{{else}}{{"Unassigned"| printf "%-10s" }}{{end}} ](fg-blue) [{{ .fields.status.name | printf "%-12s"}}](fg-blue) [{{ dateFormat "2006-01-02" .fields.created }}](fg-blue)/[{{ dateFormat "2006-01-02T15:04" .fields.updated }}](fg-green)  {{ .fields.summary | printf "%-75s"}}`
	default_event_view_template = `
issue: [{{ .key }}](fg-red)
summary: [{{ .fields.summary }}](fg-blue)

self: {{ .self }}
browse: ENDPOINT/browse/{{ .key }}
priority: {{ .fields.priority.name }}
status: {{ .fields.status.name }}
votes: {{ .fields.votes.votes }}
created: {{ .fields.created }}
updated: {{ .fields.updated }}
assignee: {{ if .fields.assignee }}{{ .fields.assignee.name }}{{end}}
reporter: {{ if .fields.reporter }}{{ .fields.reporter.name }}{{end}}
issuetype: {{ .fields.issuetype.name }}
{{if eq .fields.issuetype.name "Epic" }}epic_links: [<click here to show>](fg-red){{end}}
{{if .fields.customfield_10001 }}epic: [{{ .fields.customfield_10001 }}](fg-red){{end}}
{{if .fields.parent }}parent: [{{ .fields.parent.key }}](fg-red) -- {{ .fields.parent.fields.summary }}{{end}}
subtasks:
{{ range .fields.subtasks }}  - [{{ .key }}](fg-red)[{{.fields.status.name}}] -- {{.fields.summary}}
{{end}}

[labels:](fg-green){{ range .fields.labels }} {{ . }}{{end}}
[components:](fg-green){{ range .fields.components }} {{ .name }}{{end}}
[watchers:](fg-green){{ range .fields.customfield_10304 }} {{ .name }}{{end}}
[blockers:](fg-green)
{{ range .fields.issuelinks }}{{if .outwardIssue}}  - [{{ .outwardIssue.key }}](fg-red)[{{.outwardIssue.fields.status.name}}] -- {{.outwardIssue.fields.summary}}
{{end}}{{end}}
[depends:](fg-green)
{{ range .fields.issuelinks }}{{if .inwardIssue}}  - [{{ .inwardIssue.key }}](fg-red)[{{.inwardIssue.fields.status.name}}] -- {{.inwardIssue.fields.summary}}
{{end}}{{end}}

[description:](fg-green)

  {{ or .fields.description "" | indent 2 }}

[comments:](fg-green)

{{ range .fields.comment.comments }}
  - [{{.author.name}} at {{.created}}](fg-blue)
    {{ or .body "" | indent 4}}

{{end}}
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
