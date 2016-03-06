tessen
------

tessen is a command line tool for accessing Uchiwa/Sensu data.

It is built around the excellent [termui](https://github.com/gizak/termui) library and
the essential [JQ](https://stedolan.github.io/jq/) CLI tool.

It aims to be similar to familiar tools like vim, tig, and less.

A lot of this code is shared with
[go-jira-ui](https://github.com/mikepea/go-jira-ui), and takes much inspiration
and code from the excellent [go-jira](https://github.com/Netflix-Skunkworks/go-jira) library
(though does not depend on it)

In order to use this, you should configure at least one 'source',
with endpoint and provider, to retrieve event data from, along with at least
one query against that source:

    $ cat ~/.tessen.d/config.yml
    ---
    sources:
      - name: example-uchiwa
        provider: uchiwa
        endpoint: http://uchiwa.example.com:3001/
    queries:
      - name: All Uchiwa
        source: example-uchiwa
        filter: true

This should be all that's needed to get going.

### Installation

    # Make sure you have GOPATH and GOBIN set appropriately first:
    # eg:
    #   export GOPATH=$HOME/go
    #   export GOBIN=$GOPATH/bin
    #   mkdir -p $GOPATH
    #   export PATH=$PATH:$GOBIN
    go get -v github.com/mikepea/tessen/tessen

### Features

* Supply your own JQ boolean expressions as queries to view
* View event data with a custom template

### Usage

`tessen` is intended to mirror the options of go-jira's `jira` tool, where
useful:

    tessen             # standard usage
    tessen -h          # help page

### Basic keys

Actions:

    <enter>      - select query/event
    h            - show help page

Commands (like vim/tig/less):

    :query {JQ filter}            - display results of JQ boolean expression
    :help                          - show help page
    :<up>                          - select previous command
    :quit or :q                    - quit

Searching:

    /{regex}                       - search down
    ?{regex}                       - search up

Navigation:

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


### Configuration

It is very much recommended to read the
[go-jira](https://github.com/Netflix-Skunkworks/go-jira) documentation,
particularly surrounding the .jira.d configuration directories. tessen uses
this same mechanism, so can be used to load per-project defaults. It also
leverages the templating engine, so you can customise the view of both the
query output (use 'query_results' template), and the issue 'event_view' template.

tessen reads its own `config.yml` file in its .tessen.d directories:

    $ cat ~/.tessen.d/config.yml:
    ---
    sources:
      - name: example-uchiwa
        provider: uchiwa
        endpoint: http://uchiwa.example.com:3001/
    queries:
      - name:   'page is set'
        filter: '.query.page == true'
        source: 'example-uchiwa'
      - name:   'team is webdev and not paging'
        filter: '.query.team == "webdev" and .query.page == false'
        source: 'example-uchiwa'

Learning JQ is highly recommended. See [the JQ
manual](https://stedolan.github.io/jq/manual/) and
[tutorial](https://stedolan.github.io/jq/tutorial/), particularly around
boolean expressions. Ultimately, Tessen is doing filters using 'select'.
