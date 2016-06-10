package tessen

import (
	"fmt"
	"regexp"

	ui "github.com/gizak/termui"
)

type QueryResult struct {
	Id      string
	Summary string
	Data    interface{}
}

type QueryResultsPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	ActiveQuery Query
}

func GetQueryResults(query Query) []interface{} {
	s := FindSourceByName(query.Source.Name)
	if s == nil {
		return nil
	}
	if s.Provider == "uchiwa" {
		return GetFilteredListOfUchiwaEvents(query, &s.CachedData)
	} else if s.Provider == "pagerduty" {
		return GetFilteredListOfPagerDutyEvents(query, &s.CachedData)
	} else if s.Provider == "chronos" {
		return GetFilteredListOfChronosJobs(query, &s.CachedData)
	} else {
		log.Errorf("Unsupported provider %q", s.Provider)
		return nil
	}
}

func (p *QueryResultsPage) GetSelectedQueryResultId() string {
	qr := p.cachedResults.([]interface{})[p.selectedLine]
	return qr.(QueryResult).Id
}

func (p *QueryResultsPage) Search() {
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
		if s.re.MatchString(cr.(QueryResult).Summary) {
			p.SetSelectedLine(i)
			p.Update()
			break
		}
	}
}

func (p *QueryResultsPage) SelectItem() {
	if len(p.cachedResults.([]interface{})) == 0 {
		return
	}
	id := p.GetSelectedQueryResultId()
	if id == "" {
		return
	}
	q := new(ShowDetailPage)
	q.EventId = id
	q.Source = p.ActiveQuery.Source
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
	q.cachedResults = nil
	queryResultsPage = q
	changePage()
	q.Create()
}

func (p *QueryResultsPage) markActiveLine() {
	if p.cachedResults == nil {
		return
	}
	for i, v := range p.cachedResults.([]interface{}) {
		selected := ""
		s := v.(QueryResult).Summary
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
			if s == "" {
				s = " "
			} else if ok, _ := regexp.MatchString(`\[.+\]\((fg|bg)-[a-z]{1,6}\)`, s); ok {
				r := regexp.MustCompile(`\[(.*?)\]\((fg|bg)-[a-z]{1,6}\)`)
				s = r.ReplaceAllString(s, `$1`)
			}
			p.displayLines[i] = fmt.Sprintf("[%s](%s)", s, selected)
		} else {
			p.displayLines[i] = s
		}
	}
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
	if p.cachedResults == nil {
		p.cachedResults = GetQueryResults(p.ActiveQuery)
	}
	if p.selectedLine >= len(p.cachedResults.([]interface{})) {
		p.selectedLine = len(p.cachedResults.([]interface{})) - 1
	}
	p.displayLines = make([]string, len(p.cachedResults.([]interface{})))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("%s: %s", p.ActiveQuery.Name, p.ActiveQuery.Filter)
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
