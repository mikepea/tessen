package tessen

import (
	"fmt"
	"os"

	"github.com/coryb/optigo"
	ui "github.com/gizak/termui"
	"github.com/op/go-logging"
)

var exitNow = false

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
	log    = logging.MustGetLogger("uchiwaui")
	format = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
)

var cliOpts map[string]interface{}

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
  uchiwa-ui

General Options:
  -e --endpoint=URI   URI to use for uchiwa
  -h --help           Show this usage
  -u --user=USER      Username to use for authenticaion
  -v --verbose        Increase output logging
  --version           Print version

`)
		printer(output)
	}

	commands := map[string]string{
		"list": "list",
		"ls":   "list",
	}

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
		"u|user=s":        setopt,
		"endpoint=s":      setopt,
		"q|query=s":       setopt,
		"f|queryfields=s": setopt,
		"t|template=s":    setopt,
		"m|max_wrap=i":    setopt,
		"skip_login":      setopt,
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
	case "list":
		queryResultsPage = new(QueryResultsPage)
		if query := cliOpts["query"]; query == nil {
			log.Error("Must supply a --query option to %q", command)
			os.Exit(1)
		}
	case "toplevel":
		currentPage = queryPage
	default:
		log.Error("Unknown command %s", command)
		os.Exit(1)
	}

	for exitNow != true {

		currentPage.Create()
		ui.Loop()

	}

}
