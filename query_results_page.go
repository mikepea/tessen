package tessen

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type QueryResultsPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	ActiveQuery Query
}

func (p *QueryResultsPage) Search() {
	s := p.ActiveSearch
	n := len(p.cachedResults)
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
		if s.re.MatchString(p.cachedResults[i]) {
			p.SetSelectedLine(i)
			p.Update()
			break
		}
	}
}

func (p *QueryResultsPage) SelectItem() {
	if len(p.cachedResults) == 0 {
		return
	}
	q := new(ShowDetailPage)
	currentPage = q
	q.Create()
	changePage()
}

func (p *QueryResultsPage) GoBack() {
	currentPage = queryPage
	changePage()
}

func (p *QueryResultsPage) Update() {
	ls := p.uiList
	log.Debugf("QueryResultsPage.Update(): self:        %s (%p), ls: (%p)", p.Id(), p, ls)
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *QueryResultsPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	queryResultsPage = q
	changePage()
	q.Create()
}

func (p *QueryResultsPage) Create() {
	log.Debugf("QueryResultsPage.Create(): self:        %s (%p)", p.Id(), p)
	log.Debugf("QueryResultsPage.Create(): currentPage: %s (%p)", currentPage.Id(), currentPage)
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	query := p.ActiveQuery.Filter
	if len(p.cachedResults) == 0 {
		p.cachedResults = GetFilteredListOfEvents(query, &eventData)
	}
	if p.selectedLine >= len(p.cachedResults) {
		p.selectedLine = len(p.cachedResults) - 1
	}
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("%s: %s", p.ActiveQuery.Name, p.ActiveQuery.Filter)
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
