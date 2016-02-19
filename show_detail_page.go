package tessen

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	ui "github.com/gizak/termui"
)

const (
	defaultMaxWrapWidth = 100
)

type ShowDetailPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	EventId string
	event   interface{}
	opts    map[string]interface{}
}

func FetchEvent(id string, eventData *[]map[string]interface{}) interface{} {
	for _, ev := range *eventData {
		if ev["_id"].(string) == id {
			return ev
		}
	}
	return nil
}

func GetEventAsLines(data interface{}) []interface{} {
	buf := new(bytes.Buffer)
	results := make([]interface{}, 0)
	//template := GetTemplate("event_view")
	template := GetTemplate("debug")
	log.Debugf("GetEventAsLines: template = %q", template)
	RunTemplate(template, data, buf)
	for _, v := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		results = append(results, v)
	}
	return results
}

func (p *ShowDetailPage) Search() {
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

func (p *ShowDetailPage) Id() string {
	if p.EventId != "" {
		return fmt.Sprintf("%s", p.EventId)
	} else {
		log.Errorf("No EventId is set on %p", p)
		return ""
	}
}

func (p *ShowDetailPage) PreviousPara() {
	newDisplayLine := 0
	if p.selectedLine == 0 {
		return
	}
	for i := p.selectedLine - 1; i > 0; i-- {
		cr := p.cachedResults.([]interface{})[i]
		if ok, _ := regexp.MatchString(`^\s*$`, cr.(string)); ok {
			newDisplayLine = i
			break
		}
	}
	p.PreviousLine(p.selectedLine - newDisplayLine)
}

func (p *ShowDetailPage) NextPara() {
	newDisplayLine := len(p.cachedResults.([]interface{})) - 1
	if p.selectedLine == newDisplayLine {
		return
	}
	for i := p.selectedLine + 1; i < len(p.cachedResults.([]interface{})); i++ {
		cr := p.cachedResults.([]interface{})[i]
		if ok, _ := regexp.MatchString(`^\s*$`, cr.(string)); ok {
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
	q.cachedResults = nil
	q.event = nil
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
	ui.Clear()
	ls := ui.NewList()
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	p.uiList = ls
	if p.event == nil {
		p.event = FetchEvent(p.EventId, &eventData)
	}
	if p.cachedResults == nil {
		p.cachedResults = GetEventAsLines(p.event)
	}
	p.displayLines = make([]string, len(p.cachedResults.([]interface{})))
	if p.selectedLine >= len(p.cachedResults.([]interface{})) {
		p.selectedLine = len(p.cachedResults.([]interface{})) - 1
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