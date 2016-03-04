package tessen

import (
	"fmt"
	"os"
	"time"

	"github.com/coryb/optigo"
	ui "github.com/gizak/termui"
	"github.com/op/go-logging"
)

var exitNow = false
var defaultRefreshInterval = 30

type EditPager interface {
	DeleteRuneBackward()
	InsertRune(r rune)
	Update()
	Create()
}

type TicketCommander interface {
	ActiveTicketId() string
	Refresh()
}

type Searcher interface {
	SetSearch(string)
	Search()
}

type CommandBoxer interface {
	SetCommandMode(bool)
	ExecuteCommand()
	CommandMode() bool
	CommandBar() *CommandBar
	Update()
}

type GoBacker interface {
	GoBack()
}

type Refresher interface {
	Refresh()
}

type ItemSelecter interface {
	SelectItem()
}

type TicketEditer interface {
	EditTicket()
}

type TicketCommenter interface {
	CommentTicket()
}

type PagePager interface {
	NextLine(int)
	PreviousLine(int)
	NextPara()
	PreviousPara()
	NextPage()
	PreviousPage()
	TopOfPage()
	BottomOfPage()
	IsPopulated() bool
	Update()
}

type Navigable interface {
	Create()
	Update()
	Id() string
}

type Source struct {
	Name       string
	Provider   string
	Endpoint   string
	CachedData interface{}
}

var currentPage Navigable
var previousPage Navigable

var queryPage *QueryPage
var helpPage *HelpPage
var queryResultsPage *QueryResultsPage
var commandBar *CommandBar

func changePage() {
	switch currentPage.(type) {
	case *QueryPage:
		log.Debugf("changePage: QueryPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *QueryResultsPage:
		log.Debugf("changePage: QueryResultsPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *ShowDetailPage:
		log.Debugf("changePage: ShowDetailPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *HelpPage:
		log.Debugf("changePage: HelpPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	}
}

var (
	log    = logging.MustGetLogger("tessen")
	format = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
)

var cliOpts map[string]interface{}
var sources []*Source

func getSources(opts map[string]interface{}) []*Source {
	if _, ok := opts["sources"]; !ok {
		log.Fatal("No sources specified, exiting")
	}
	sources = make([]*Source, 0)
	for _, s := range opts["sources"].([]interface{}) {
		source := s.(map[interface{}]interface{})
		name := source["name"].(string)
		endpoint := source["endpoint"].(string)
		provider := source["provider"].(string)
		sources = append(sources, &Source{name, provider, endpoint, nil})
	}
	return sources
}

func collectSource(s *Source, seconds int) {
	timer := time.NewTimer(time.Duration(seconds) * time.Second)
	<-timer.C
	var err error
	if s.Provider == "uchiwa" {
		log.Debugf("Collecting uchiwa data")
		s.CachedData, err = FetchUchiwaEvents(s.Endpoint)
	} else if s.Provider == "pagerduty" {
		log.Debugf("Collecting pagerduty data")
		s.CachedData, err = FetchPagerDutyEvents(s.Endpoint)
	} else {
		log.Errorf("Cannot collect from source %q, unimplemented backend type %q", s.Name, s.Provider)
	}
	if err != nil {
		log.Errorf("Error fetching source data for %q: %q\n", s.Name, err)
	}
}

func FindSourceByName(name string) *Source {
	for _, s := range sources {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func Run() {

	var err error
	logging.SetLevel(logging.NOTICE, "")

	usage := func(ok bool) {
		printer := fmt.Printf
		if !ok {
			printer = func(format string, args ...interface{}) (int, error) {
				return fmt.Fprintf(os.Stderr, format, args...)
			}
			defer func() {
				os.Exit(1)
			}()
		} else {
			defer func() {
				os.Exit(0)
			}()
		}
		output := fmt.Sprintf(`
Usage:
  tessen

General Options:
  -h --help           Show this usage
  -v --verbose        Increase output logging
  --version           Print version

`)
		printer(output)
	}

	commands := map[string]string{}

	cliOpts = make(map[string]interface{})
	setopt := func(name string, value interface{}) {
		cliOpts[name] = value
	}

	op := optigo.NewDirectAssignParser(map[string]interface{}{
		"h|help": usage,
		"version": func() {
			fmt.Println(fmt.Sprintf("version: %s", VERSION))
			os.Exit(0)
		},
		"v|verbose+": func() {
			logging.SetLevel(logging.GetLevel("")+1, "")
		},
		"l|listen=s": setopt,
		"noui":       setopt,
	})

	if err := op.ProcessAll(os.Args[1:]); err != nil {
		log.Error("%s", err)
		usage(false)
	}
	args := op.Args

	var command string
	if len(args) > 0 {
		if alias, ok := commands[args[0]]; ok {
			command = alias
			args = args[1:]
		} else {
			command = "view"
			args = args[0:]
		}
	} else {
		command = "toplevel"
	}

	opts := getOpts()
	sources := getSources(opts)

	if _, ok := opts["noui"]; !ok {
		err = ui.Init()
		if err != nil {
			panic(err)
		}
		defer ui.Close()

		registerKeyboardHandlers()

		queryPage = new(QueryPage)
		helpPage = new(HelpPage)
		commandBar = new(CommandBar)

		switch command {
		case "toplevel":
			currentPage = queryPage
		default:
			log.Error("Unknown command %s", command)
			os.Exit(1)
		}
	}

	for _, s := range sources {
		collectSource(s, 0)
	}

	for _, s := range sources {
		go func() {
			for {
				collectSource(s, defaultRefreshInterval)
				// TODO: we need a way of calling this only if the data has changed, or if the
				//  'currentPage' relates to this dataset.
				if obj, ok := currentPage.(Refresher); ok {
					obj.Refresh()
				}
			}
		}()
	}

	if l, ok := opts["listen"]; ok {
		log.Debugf("Starting http dashboard on %s", l.(string))
		go func() {
			log.Fatal(StartHttpDashboard(l.(string)))
		}()
	}

	for exitNow != true {

		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}

		if _, ok := opts["noui"]; !ok {
			currentPage.Create()
			ui.Loop()
		}
		log.Debug("End of exitNow loop")

	}

	log.Debug("Normal exit, woohoo!")

}
