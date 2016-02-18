package tessen

import (
	"strings"
)

type CommandBarFragment struct {
	commandBar  *CommandBar
	commandMode bool
}

func (p *CommandBarFragment) ExecuteCommand() {
	command := string(p.commandBar.text)
	if command == "" {
		return
	}
	commandMode := string([]rune(command)[0])
	switch commandMode {
	case "/":
		log.Debugf("Search down: %q", command)
		if obj, ok := currentPage.(Searcher); ok {
			obj.SetSearch(command)
			obj.Search()
		}
	case "?":
		log.Debugf("Search up: %q", command)
		if obj, ok := currentPage.(Searcher); ok {
			obj.SetSearch(command)
			obj.Search()
		}
	case ":":
		log.Debugf("Command: %q", command)
		handleCommand(command)
	}
}

func handleCommand(command string) {
	if len(command) < 2 {
		// must be :something
		return
	}
	fields := strings.Fields(string(command[1:]))
	action := fields[0]
	var args []string
	if len(fields) > 1 {
		args = fields[1:]
	}
	log.Debugf("handleCommand: action %q, args %s", action, args)
	switch {
	case action == "q" || action == "quit":
		handleQuit()
	case action == "help":
		handleHelp()
	case action == "query":
		n := len(":query ")
		if len(command) > n {
			handleQueryCommand(string(command[(n - 1):]))
		}
	}
}

func handleQueryCommand(query string) {
	log.Debugf("handleQueryCommand: query %q", query)
	if query == "" {
		return
	}
	q := new(QueryResultsPage)
	q.ActiveQuery.Name = "adhoc query"
	q.ActiveQuery.Filter = query
	currentPage = q
	changePage()
}

func (p *CommandBarFragment) SetCommandMode(mode bool) {
	p.commandMode = mode
}

func (p *CommandBarFragment) CommandMode() bool {
	return p.commandMode
}

func (p *CommandBarFragment) CommandBar() *CommandBar {
	return p.commandBar
}
