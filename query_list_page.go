package tessen

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type Query struct {
	Name     string
	Filter   string
	Template string
}

type QueryPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	cachedResults []Query
}

var baseQueries = []Query{
	Query{"All Events", "true", ""},
	Query{"---", "", ""}, // no-op line in UI
}

func getQueries() (queries []Query) {
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
			queries = append(queries, Query{q2["name"], q2["filter"], q2["template"]})
		}
	}
	return append(baseQueries, queries...)
}

func (p *QueryPage) Search() {
	s := p.ActiveSearch
	log.Debugf("QueryPage: search! %q", s.command)
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
		if s.re.MatchString(p.cachedResults[i].Name) {
			log.Debugf("Match found, line %d", i)
			p.SetSelectedLine(i)
			p.Update()
			break
		}
	}
}

func (p *QueryPage) IsPopulated() bool {
	if len(p.cachedResults) > 0 {
		return true
	} else {
		return false
	}
}

func (p *QueryPage) SetSelectedLine(line int) {
	if line > 0 && line < len(p.cachedResults) {
		p.selectedLine = line
		p.FixFirstDisplayLine(0)
	}
}

func (p *QueryPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
			p.displayLines[i] = fmt.Sprintf("[%-50s | %s](%s)", v.Name, v.Filter, selected)
		} else {
			p.displayLines[i] = fmt.Sprintf("%-50s [|](fg-blue) [%s](fg-green)", v.Name, v.Filter)
		}
	}
}

func (p *QueryPage) PreviousPara() {
	newDisplayLine := 0
	if p.selectedLine == 0 {
		return
	}
	for i := p.selectedLine - 1; i > 0; i-- {
		if p.cachedResults[i].Filter == "" {
			newDisplayLine = i
			break
		}
	}
	p.PreviousLine(p.selectedLine - newDisplayLine)
}

func (p *QueryPage) NextPara() {
	newDisplayLine := len(p.cachedResults) - 1
	if p.selectedLine == newDisplayLine {
		return
	}
	for i := p.selectedLine + 1; i < len(p.cachedResults); i++ {
		if p.cachedResults[i].Filter == "" {
			newDisplayLine = i
			break
		}
	}
	p.NextLine(newDisplayLine - p.selectedLine)
}

func (p *QueryPage) BottomOfPage() {
	p.selectedLine = len(p.cachedResults) - 1
	firstLine := p.selectedLine - (p.uiList.Height - 3)
	if firstLine > 0 {
		p.firstDisplayLine = firstLine
	} else {
		p.firstDisplayLine = 0
	}
}

func (p *QueryPage) SelectedQuery() Query {
	return p.cachedResults[p.selectedLine]
}

func (p *QueryPage) SelectItem() {
	if p.SelectedQuery().Filter == "" {
		return
	}
	q := new(QueryResultsPage)
	q.ActiveQuery = p.SelectedQuery()
	queryResultsPage = q
	currentPage = queryResultsPage
	changePage()
}

func (p *QueryPage) Update() {
	ls := p.uiList
	log.Debugf("QueryPage.Update(): self:        %s (%p), ls: (%p)", p.Id(), p, ls)
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *QueryPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]Query, 0)
	changePage()
	q.Create()
}

func (p *QueryPage) Create() {
	log.Debugf("QueryPage.Create(): self:        %s (%p)", p.Id(), p)
	log.Debugf("QueryPage.Create(): currentPage: %s (%p)", currentPage.Id(), currentPage)
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	p.cachedResults = getQueries()
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Queries"
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
