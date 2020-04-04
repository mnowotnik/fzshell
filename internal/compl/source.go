package compl

import (
	"os"
	"os/exec"
	"strings"

	"text/template"

	"github.com/pkg/errors"
)

type CompletionSource struct {
	Command       string `yaml:"cmd"`
	ItemTmpl      string `yaml:"map"`
	ViewTmpl      string `yaml:"mapView"`
	PreviewTmpl   string `yaml:"preview"`
	FilterTmpl    string `yaml:"filter"`
	LineSeparator string `yaml:"lineSeparator"`
	HeaderLines   int    `yaml:"headerLines"`
	PreferHeader  bool   `yaml:"selectHeader"`
}

type CompletionEntry struct {
	Item    string
	View    string
	Line    string
	Preview *Preview
}

func (cs *CompletionSource) generateEntries(match *MatchResult) (CompletionResult, error) {
	var items []string

	out, err := executeInShell(cs.Command, match.subExp, match.subExpNamed)

	if err != nil {
		return CompletionResult{}, errors.New(err.Error() + string(out))
	}
	if cs.LineSeparator == "" {
		cs.LineSeparator = "\n"
	}
	items = strings.Split(out, cs.LineSeparator)
	header := items[:cs.HeaderLines]
	items = items[cs.HeaderLines:]

	entries := []CompletionEntry{}
	viewTmpl, err := createTemplate(cs.ViewTmpl)
	if err != nil {
		return CompletionResult{}, err
	}

	itemTmpl, err := createTemplate(cs.ItemTmpl)
	if err != nil {
		return CompletionResult{}, err
	}

	filterTmpl, err := createTemplate(cs.FilterTmpl)
	if err != nil {
		return CompletionResult{}, err
	}

	previewTmpl, err := createTemplate(cs.PreviewTmpl)
	if err != nil {
		return CompletionResult{}, err
	}

	tmplData := make(map[string]string)
	for k, v := range match.subExpNamed {
		tmplData[k] = v
	}
	for _, item := range items {
		tmplData["item"] = item
		if item == "" {
			continue
		}
		if filterTmpl != nil {
			if buf, err := renderFromTemplate(filterTmpl, match.subExp, tmplData); err != nil {
				return CompletionResult{}, err
			} else if buf.String() == "false" {
				continue
			}
		}

		var itemStr string

		if cs.ItemTmpl != "" {
			if buf, err := renderFromTemplate(itemTmpl, match.subExp, tmplData); err != nil {
				return CompletionResult{}, err
			} else {
				itemStr = buf.String()
			}
		} else {
			itemStr = item
		}

		var viewStr string
		if cs.ViewTmpl != "" {
			if buf, err := renderFromTemplate(viewTmpl, match.subExp, tmplData); err != nil {
				return CompletionResult{}, err
			} else {
				viewStr = buf.String()
			}
		} else {
			viewStr = item
		}
		var preview *Preview
		if previewTmpl != nil {
			preview = &Preview{PreviewTmpl: previewTmpl, MatchResult: match}
		}
		entries = append(entries, CompletionEntry{Item: itemStr, View: viewStr, Line: item, Preview: preview})
	}
	return CompletionResult{Header: header, Entries: entries}, nil
}

func createTemplate(tmplStr string) (*template.Template, error) {
	if tmplStr != "" {
		tmpl, err := parseTemplate(tmplStr)
		if err != nil {
			return tmpl, errors.Wrap(err, "Error parsing template: "+tmplStr)
		}
		return tmpl, nil
	}
	return nil, nil
}

func executeInShell(expr string, args []string, kwargs map[string]string) (string, error) {
	expr = "set -o pipefail;" + expr
	cArgs := []string{"-c", expr}
	cArgs = append(cArgs, "/usr/bin/bash")
	cArgs = append(cArgs, args...)
	cmd := exec.Command("bash", cArgs...)
	env := os.Environ()
	for k, v := range kwargs {
		env = append(env, k+"="+v)
	}
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	return string(out), err
}
