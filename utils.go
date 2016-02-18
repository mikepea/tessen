package tessen

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Netflix-Skunkworks/go-jira"
	ui "github.com/gizak/termui"
	"github.com/mitchellh/go-wordwrap"
	"gopkg.in/coryb/yaml.v2"
)

func HelpTextAsStrings(data interface{}, templateName string) []string {
	return strings.Split("", "\n")
}

func WrapText(lines []string, maxWidth uint) []string {
	out := make([]string, 0)
	insideNoformatBlock := false
	insideCodeBlock := false
	for _, line := range lines {
		if matched, _ := regexp.MatchString(`^\s+\{code`, line); matched {
			insideCodeBlock = !insideCodeBlock
		} else if strings.TrimSpace(line) == "{noformat}" {
			insideNoformatBlock = !insideNoformatBlock
		}
		if maxWidth == 0 || uint(len(line)) < maxWidth || insideCodeBlock || insideNoformatBlock {
			out = append(out, line)
			continue
		}
		if matched, _ := regexp.MatchString(`^[a-z_]+:\s`, line); matched {
			// don't futz with single line field+value.
			// If they are too long, that's their fault.
			out = append(out, line)
			continue
		}
		// wrap text, but preserve indenting
		re := regexp.MustCompile(`^\s*`)
		indenting := re.FindString(line)
		wrappedLines := strings.Split(wordwrap.WrapString(line, maxWidth-uint(len(indenting))), "\n")
		indentedWrappedLines := make([]string, len(wrappedLines))
		for i, wl := range wrappedLines {
			if i == 0 {
				// first line already has the indent
				indentedWrappedLines[i] = wl
			} else {
				indentedWrappedLines[i] = indenting + wl
			}
		}
		out = append(out, indentedWrappedLines...)
	}
	return out
}

func parseYaml(file string, v map[string]interface{}) {
	if fh, err := ioutil.ReadFile(file); err == nil {
		log.Debugf("Parsing YAML file: %s", file)
		yaml.Unmarshal(fh, &v)
	}
}

func loadConfigs(opts map[string]interface{}) {
	paths := jira.FindParentPaths(".uchiwa-ui.d/config.yml")
	paths = append([]string{"/etc/go-uchiwa-ui.yml"}, paths...)

	// iterate paths in reverse
	for i := len(paths) - 1; i >= 0; i-- {
		file := paths[i]
		if _, err := os.Stat(file); err == nil {
			tmp := make(map[string]interface{})
			parseYaml(file, tmp)
			for k, v := range tmp {
				if _, ok := opts[k]; !ok {
					log.Debugf("Setting %q to %#v from %s", k, v, file)
					opts[k] = v
				}
			}
		}
	}
}

func getOpts() map[string]interface{} {
	user := os.Getenv("USER")
	home := os.Getenv("HOME")

	opts := make(map[string]interface{})
	defaults := map[string]interface{}{
		"user":      user,
		"endpoint":  os.Getenv("JIRA_ENDPOINT"),
		"directory": fmt.Sprintf("%s/.jira.d/templates", home),
		"quiet":     true,
	}

	for k, v := range cliOpts {
		if _, ok := opts[k]; !ok {
			log.Debugf("Setting %q to %#v from cli options", k, v)
			opts[k] = v
		}
	}

	loadConfigs(opts)
	for k, v := range defaults {
		if _, ok := opts[k]; !ok {
			log.Debugf("Setting %q to %#v from defaults", k, v)
			opts[k] = v
		}
	}
	return opts
}

func lastLineDisplayed(ls *ui.List, firstLine int, correction int) int {
	return firstLine + ls.Height - correction
}
