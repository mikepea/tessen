package tessen

import (
	"bytes"
	"strings"

	ui "github.com/gizak/termui"
)

type HelpPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
}

func HelpTextAsStrings(data interface{}, templateName string) []interface{} {
	buf := new(bytes.Buffer)
	results := make([]interface{}, 0)
	template := GetTemplate(templateName)
	log.Debugf("HelpTextAsStrings: template = %q", template)
	RunTemplate(template, data, buf)
	for _, v := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		results = append(results, v)
	}
	return results
}

func (p *HelpPage) Search() {
	s := p.ActiveSearch
	n := len(p.cachedResults.([]interface{}))
	if s.command == "" {
		return
	}
	increment := 1
	if s.directionUp {
		increment = -1
	}
	// we use modulo here so we can loop through every line.
	// adding 'n' means we never have '-1 % n'.
	startLine := (p.selectedLine + n + increment) % n
	for i := startLine; i != p.selectedLine; i = (i + increment + n) % n {
		cr := p.cachedResults.([]interface{})[i]
		if s.re.MatchString(cr.(string)) {
			p.SetSelectedLine(i)
			p.Update()
			break
		}
	}
}

func (p *HelpPage) GoBack() {
	currentPage = previousPage
	changePage()
}

func (p *HelpPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *HelpPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	if p.cachedResults == nil {
		p.cachedResults = HelpTextAsStrings(nil, "help")
	}
	p.displayLines = make([]string, len(p.cachedResults.([]interface{})))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Help"
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
