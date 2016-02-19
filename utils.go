package tessen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	ui "github.com/gizak/termui"
	"gopkg.in/coryb/yaml.v2"
)

func TemplatedEvent(id string, templateName string, eventData interface{}) []string {
	if templateName == "" {
		templateName = "event_view"
	}
	template := GetTemplate(templateName)
	buf := new(bytes.Buffer)
	if template == "" {
		template = default_list_template
	}
	RunTemplate(template, eventData, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func parseYaml(file string, v map[string]interface{}) {
	if fh, err := ioutil.ReadFile(file); err == nil {
		log.Debugf("Parsing YAML file: %s", file)
		yaml.Unmarshal(fh, &v)
	}
}

func FindParentPaths(fileName string) []string {
	cwd, _ := os.Getwd()

	paths := make([]string, 0)

	// special case if homedir is not in current path then check there anyway
	homedir := os.Getenv("HOME")
	if !strings.HasPrefix(cwd, homedir) {
		file := fmt.Sprintf("%s/%s", homedir, fileName)
		if _, err := os.Stat(file); err == nil {
			paths = append(paths, file)
		}
	}

	var dir string
	for _, part := range strings.Split(cwd, string(os.PathSeparator)) {
		if dir == "/" {
			dir = fmt.Sprintf("/%s", part)
		} else {
			dir = fmt.Sprintf("%s/%s", dir, part)
		}
		file := fmt.Sprintf("%s/%s", dir, fileName)
		if _, err := os.Stat(file); err == nil {
			paths = append(paths, file)
		}
	}
	return paths
}

func loadConfigs(opts map[string]interface{}) {
	paths := FindParentPaths(".tessen.d/config.yml")
	paths = append([]string{"/etc/tessen.yml"}, paths...)

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
	home := os.Getenv("HOME")

	opts := make(map[string]interface{})
	defaults := map[string]interface{}{
		"endpoint":  os.Getenv("UCHIWA_ENDPOINT"),
		"directory": fmt.Sprintf("%s/.tessen.d/templates", home),
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
