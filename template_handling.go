package tessen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func readFile(file string) string {
	var bytes []byte
	var err error
	if bytes, err = ioutil.ReadFile(file); err != nil {
		log.Error("Failed to read file %s: %s", file, err)
		os.Exit(1)
	}
	return string(bytes)
}

func FindClosestParentPath(fileName string) (string, error) {
	paths := FindParentPaths(fileName)
	if len(paths) > 0 {
		return paths[len(paths)-1], nil
	}
	return "", errors.New(fmt.Sprintf("%s not found in parent directory hierarchy", fileName))
}

func GetTemplate(name string) string {
	opts := getOpts()
	pathSpec := ".tessen.d/templates/%s"
	if override, ok := opts["template"].(string); ok {
		if _, err := os.Stat(override); err == nil {
			return readFile(override)
		} else {
			if file, err := FindClosestParentPath(fmt.Sprintf(pathSpec, override)); err == nil {
				return readFile(file)
			}
			if dflt, ok := all_templates[override]; ok {
				return dflt
			}
		}
	}
	if file, err := FindClosestParentPath(fmt.Sprintf(pathSpec, name)); err != nil {
		return all_templates[name]
	} else {
		return readFile(file)
	}
}

func dateFormat(format string, content string) (string, error) {
	if t, err := time.Parse("2006-01-02T15:04:05.000-0700", content); err != nil {
		return "", err
	} else {
		return t.Format(format), nil
	}
}

func colorizedSensuStatus(statusCode int) (string, error) {
	status := ""
	switch statusCode {
	case 0:
		status = "[OK  ](fg-green)"
	case 1:
		status = "[WARN](fg-yellow)"
	case 2:
		status = "[CRIT](fg-red)"
	case 3:
		status = "[UNKN](fg-blue)"
	default:
		return "", errors.New(fmt.Sprintf("Invalid statusCode %q", statusCode))
	}
	return status, nil
}

func RunTemplate(templateContent string, data interface{}, out io.Writer) error {
	if out == nil {
		out = os.Stdout
	}

	funcs := map[string]interface{}{
		"toJson": func(content interface{}) (string, error) {
			if bytes, err := json.MarshalIndent(content, "", "    "); err != nil {
				return "", err
			} else {
				return string(bytes), nil
			}
		},
		"append": func(more string, content interface{}) (string, error) {
			switch value := content.(type) {
			case string:
				return string(append([]byte(content.(string)), []byte(more)...)), nil
			case []byte:
				return string(append(content.([]byte), []byte(more)...)), nil
			default:
				return "", errors.New(fmt.Sprintf("Unknown type: %s", value))
			}
		},
		"indent": func(spaces int, content string) string {
			indent := make([]rune, spaces+1, spaces+1)
			indent[0] = '\n'
			for i := 1; i < spaces+1; i += 1 {
				indent[i] = ' '
			}

			lineSeps := []rune{'\n', '\u0085', '\u2028', '\u2029'}
			for _, sep := range lineSeps {
				indent[0] = sep
				content = strings.Replace(content, string(sep), string(indent), -1)
			}
			return content

		},
		"comment": func(content string) string {
			lineSeps := []rune{'\n', '\u0085', '\u2028', '\u2029'}
			for _, sep := range lineSeps {
				content = strings.Replace(content, string(sep), string([]rune{sep, '#', ' '}), -1)
			}
			return content
		},
		"split": func(sep string, content string) []string {
			return strings.Split(content, sep)
		},
		"join": func(sep string, content []interface{}) string {
			vals := make([]string, len(content))
			for i, v := range content {
				vals[i] = v.(string)
			}
			return strings.Join(vals, sep)
		},
		"abbrev": func(max int, content string) string {
			if len(content) > max {
				var buffer bytes.Buffer
				buffer.WriteString(content[:max-3])
				buffer.WriteString("...")
				return buffer.String()
			}
			return content
		},
		"rep": func(count int, content string) string {
			var buffer bytes.Buffer
			for i := 0; i < count; i += 1 {
				buffer.WriteString(content)
			}
			return buffer.String()
		},
		"dateFormat": func(format string, content string) (string, error) {
			return dateFormat(format, content)
		},
		"colorizedSensuStatus": func(statusCode interface{}) (string, error) {
			switch statusCode := statusCode.(type) {
			case float64:
				return colorizedSensuStatus(int(statusCode))
			case string:
				if i, err := strconv.Atoi(statusCode); err == nil {
					return colorizedSensuStatus(i)
				}
			}
			return "", errors.New("colorizedSensuStatus: bad type")
		},
	}
	if tmpl, err := template.New("template").Funcs(funcs).Parse(templateContent); err != nil {
		log.Error("Failed to parse template: %s", err)
		return err
	} else {
		if err := tmpl.Execute(out, data); err != nil {
			log.Error("Failed to execute template: %s", err)
			return err
		}
	}
	return nil
}
