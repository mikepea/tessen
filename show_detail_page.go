package tessen

import (
	"fmt"
	"regexp"

	ui "github.com/gizak/termui"
)

const (
	defaultMaxWrapWidth = 100
)

type ShowDetailPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	MaxWrapWidth uint
	Template     string
	apiBody      interface{}
	WrapWidth    uint
	opts         map[string]interface{}
}

func (p *ShowDetailPage) Search() {
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

func (p *ShowDetailPage) SelectItem() {
	//selected := p.cachedResults[p.selectedLine]
	return
	/*
			q := new(ShowDetailPage)
			q.TicketId = newTicketId
			currentPage = q
		changePage()
	*/
}

func (p *ShowDetailPage) Id() string {
	return "TODO"
}

func (p *ShowDetailPage) PreviousPara() {
	newDisplayLine := 0
	if p.selectedLine == 0 {
		return
	}
	for i := p.selectedLine - 1; i > 0; i-- {
		if ok, _ := regexp.MatchString(`^\s*$`, p.cachedResults[i]); ok {
			newDisplayLine = i
			break
		}
	}
	p.PreviousLine(p.selectedLine - newDisplayLine)
}

func (p *ShowDetailPage) NextPara() {
	newDisplayLine := len(p.cachedResults) - 1
	if p.selectedLine == newDisplayLine {
		return
	}
	for i := p.selectedLine + 1; i < len(p.cachedResults); i++ {
		if ok, _ := regexp.MatchString(`^\s*$`, p.cachedResults[i]); ok {
			newDisplayLine = i
			break
		}
	}
	p.NextLine(newDisplayLine - p.selectedLine)
}

func (p *ShowDetailPage) GoBack() {
	if queryResultsPage != nil {
		currentPage = queryResultsPage
	} else {
		currentPage = queryPage
	}
	changePage()
}

func (p *ShowDetailPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	q.apiBody = nil
	currentPage = q
	changePage()
	q.Create()
}

func (p *ShowDetailPage) Update() {
	ls := p.uiList
	log.Debugf("ShowDetailPage.Update(): self:        %s (%p), ls: (%p)", p.Id(), p, ls)
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *ShowDetailPage) Create() {
	log.Debugf("ShowDetailPage.Create(): self:        %s (%p)", p.Id(), p)
	log.Debugf("ShowDetailPage.Create(): currentPage: %s (%p)", currentPage.Id(), currentPage)
	p.opts = getOpts()
	/*
		if p.TicketId == "" {
			p.TicketId = ticketListPage.GetSelectedTicketId()
		}
	*/
	if p.MaxWrapWidth == 0 {
		if m := p.opts["max_wrap"]; m != nil {
			p.MaxWrapWidth = uint(m.(int64))
		} else {
			p.MaxWrapWidth = defaultMaxWrapWidth
		}
	}
	ui.Clear()
	ls := ui.NewList()
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	p.uiList = ls
	innerWidth := uint(ui.TermWidth()) - 3
	if innerWidth < p.MaxWrapWidth {
		p.WrapWidth = innerWidth
	} else {
		p.WrapWidth = p.MaxWrapWidth
	}
	if p.apiBody == nil {
		/*
			p.apiBody, _ = FetchJiraTicket(p.TicketId)
		*/
	}
	/*
		p.cachedResults = WrapText(JiraTicketAsStrings(p.apiBody, p.Template), p.WrapWidth)
	*/
	p.displayLines = make([]string, len(p.cachedResults))
	if p.selectedLine >= len(p.cachedResults) {
		p.selectedLine = len(p.cachedResults) - 1
	}
	ls.ItemFgColor = ui.ColorYellow
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Border = true
	ls.BorderLabel = fmt.Sprintf("%s", p.Id)
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
