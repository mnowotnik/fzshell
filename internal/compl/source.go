package compl

import (
	"os"
	"os/exec"

	"text/template"

	fzf "github.com/mnowotnik/fzf/src"
	"github.com/mnowotnik/fzshell/internal/utils"
	"github.com/pkg/errors"
)

type CompletionSource struct {
	Command      string      `yaml:"cmd"`
	ItemTmpl     string      `yaml:"map"`
	PreviewTmpl  string      `yaml:"preview"`
	HeaderLines  int         `yaml:"headerLines"`
	Header       interface{} `yaml:"header"`
	SelectFirst  bool        `yaml:"selectFirst"`
	PreferHeader bool        `yaml:"selectHeader"`
}

func (cs *CompletionSource) generateEntries(match *MatchResult, returnAll bool) ([]string, error) {
	if cs.Command == "" {
		return []string{}, nil
	}
	results, err := cs.pipeCommandToFzf(match.subExp, match.subExpNamed, returnAll)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return results, nil
	}

	itemTmpl, err := createTemplate(cs.ItemTmpl)
	if err != nil {
		return nil, err
	}

	tmplData := utils.CloneMap(match.subExpNamed)
	items := []string{}
	for _, item := range results {
		tmplData["item"] = item
		if item == "" {
			continue
		}

		var itemStr string

		if cs.ItemTmpl != "" {
			if buf, err := renderFromTemplate(itemTmpl, match.subExp, tmplData); err != nil {
				return nil, err
			} else {
				itemStr = buf.String()
			}
		} else {
			itemStr = item
		}

		items = append(items, itemStr)
	}
	return items, nil
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

func (cs *CompletionSource) pipeCommandToFzf(args []string, kwargs map[string]string, returnAll bool) ([]string, error) {
	os.Setenv("FZF_DEFAULT_OPTS", "")
	os.Args = []string{os.Args[0]}
	options := fzf.ParseOptions()
	if cs.PreviewTmpl != "" {
		options.Preview.Command = cs.PreviewTmpl

		previewTmpl, err := createTemplate(cs.PreviewTmpl)
		if err != nil {
			return nil, err
		}

		options.Preview.CommandGenerator = func(item string, query string) (string, error) {
			kwargsC := utils.CloneMap(kwargs)
			kwargsC["item"] = item
			buf, err := renderFromTemplate(previewTmpl, args, kwargsC)
			if err != nil {
				return "", err
			}
			return buf.String(), nil
		}
		options.HeaderLines = cs.HeaderLines
		if cs.Header != nil {
			switch v := cs.Header.(type) {
			case string:
				options.Header = []string{v}
			case []string:
				options.Header = v
			default:
				return nil, errors.New("wrong value for 'header' key")
			}
		}
		options.Select1 = cs.SelectFirst
	}
	if returnAll {
		var filter string = ""
		options.Filter = &filter
	}
	expr := "set -o pipefail;" + cs.Command
	cArgs := []string{"-c", expr}
	cArgs = append(cArgs, "/usr/bin/bash")
	cArgs = append(cArgs, args...)
	cmd := exec.Command("bash", cArgs...)
	env := os.Environ()
	for k, v := range kwargs {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	out, err := cmd.StdoutPipe()
	results := []string{}

	options.Printer = func(i string) {
		results = append(results, i)
	}
	cmd.Start()
	fzf.Run(options, out)
	return results, err
}
