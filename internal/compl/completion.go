package compl

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/mnowotnik/fzshell/internal/output"
	"github.com/mnowotnik/go-fuzzyfinder/fuzzyfinder"
)

type Completion struct {
	CompletionSource `yaml:",inline"`
	Pattern          string             `yaml:"pattern"`
	Replacement      string             `yaml:"replacement"`
	SelectFirst      bool               `yaml:"selectFirst"`
	ItemSeparator    string             `default:" " yaml:"itemSeparator"`
	Layout           string             `yaml:"layout"`
	Sources          []CompletionSource `yaml:"sources"`
}

type CompletionResult struct {
	Header  []string
	Entries []CompletionEntry
}

type CompletionOptions struct {
	ReturnAll bool
}

type ComplResult int

const (
	Matched ComplResult = iota
	NotMatched
	Aborted
	Errored
)

func NewCompletion(pattern string) *Completion {
	return &Completion{Pattern: pattern}
}

func (c *Completion) MatchAndFind(line string, options CompletionOptions) ([]string, ComplResult) {
	err := c.validate()
	if err != nil {
		output.Log().Println(err)
		return []string{}, Errored
	}
	matcher := matcher{pattern: c.Pattern}
	m, err := matcher.match(line)
	if err != nil {
		output.Log().Println(err)
		return []string{}, Errored
	}
	if m == nil {
		return nil, NotMatched
	}
	return c.execute(m, options)
}

func (c *Completion) validate() error {
	if c.Pattern == "" {
		return fmt.Errorf("pattern not set: %#v", c)
	}

	if len(c.Sources) > 0 {
		for _, s := range c.Sources {
			if s.Command == "" {
				return fmt.Errorf("source command not set: %#v", c)
			}
		}
	}
	return nil
}

func (c *Completion) execute(matchResult *MatchResult, options CompletionOptions) ([]string, ComplResult) {
	result, err := c.finder(matchResult, options.ReturnAll)
	if err != nil {
		if strings.Contains(err.Error(), "abort") {
			return []string{}, Aborted
		}
		output.Log().Println(err)
		return []string{}, Errored
	}
	if options.ReturnAll {
		return result, Matched
	}
	if str, err := c.renderLBuffer(matchResult, result); err == nil {
		return []string{str}, Matched
	} else {
		output.Log().Println(err)
		return []string{}, Errored
	}
}

func (c *Completion) generate(match *MatchResult) (CompletionResult, error) {

	if len(c.Sources) == 0 {
		if c.Command == "" {
			return CompletionResult{}, nil
		}
		return c.generateEntries(match)
	}

	result := CompletionResult{}
	for _, source := range c.Sources {
		sourceResult, error := source.generateEntries(match)
		if error != nil {
			return CompletionResult{}, error
		}
		if source.PreferHeader && len(sourceResult.Header) > 0 {
			result.Header = sourceResult.Header
		}
		result.Entries = append(result.Entries, sourceResult.Entries...)
	}
	return result, nil
}

func (c *Completion) finder(match *MatchResult, returnAll bool) ([]string, error) {

	// TODO add sticky header to fuzzyfinder
	complResult, err := c.generate(match)
	if err != nil {
		return nil, err
	}
	if len(complResult.Entries) == 0 {
		return []string{}, nil
	}

	if len(complResult.Entries) == 1 || c.SelectFirst {
		return []string{complResult.Entries[0].Item}, nil
	}
	if returnAll {
		return pie.Map(complResult.Entries, func(e CompletionEntry) string { return e.Item }), nil
	}
	layoutOpt := fuzzyfinder.WithLayout(fuzzyfinder.SmartLayout)
	if c.Layout == "concise" {
		layoutOpt = fuzzyfinder.WithLayout(fuzzyfinder.NoLayout)
	}
	preview := fuzzyfinder.WithPreviewWindow(func(i int, width int, height int) string {
		r, err := complResult.Entries[i].RenderPreview()
		if err != nil {
			output.Log().Printf("Error rendering preview: %s", err.Error())
			os.Exit(1)
		}
		return r
	})
	indices, err := fuzzyfinder.FindMulti(
		complResult.Entries,
		func(i int) fuzzyfinder.Item {
			return fuzzyfinder.Item{Value: complResult.Entries[i].Item, View: complResult.Entries[i].View}
		}, preview, layoutOpt)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, i := range indices {
		result = append(result, complResult.Entries[i].Item)
	}
	return result, nil
}

var spaces = regexp.MustCompile(" +")

func (c *Completion) renderLBuffer(matchR *MatchResult, items []string) (string, error) {
	var result string
	if c.Replacement == "" {
		result = matchR.line + strings.Join(items, c.ItemSeparator)
	} else {
		vars := make(map[string]string)
		vars["item"] = strings.Join(items, c.ItemSeparator)
		for k, v := range matchR.subExpNamed {
			vars[k] = v
		}
		buf, err := render(c.Replacement, matchR.subExp, vars)
		if err != nil {
			return "", err
		}
		result = matchR.line[:matchR.startIdx] + buf.String()
	}

	return result, nil
}
