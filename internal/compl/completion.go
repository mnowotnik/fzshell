package compl

import (
	"fmt"
	"strings"

	"github.com/mnowotnik/fzshell/internal/output"
)

type Completion struct {
	CompletionSource `yaml:",inline"`
	Pattern          string `yaml:"pattern"`
	Replacement      string `yaml:"replacement"`
	ItemSeparator    string `default:" " yaml:"itemSeparator"`
	Layout           string `yaml:"layout"`
}

type CompletionResult struct {
	Items []string
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
	return nil
}

func (c *Completion) execute(matchResult *MatchResult, options CompletionOptions) ([]string, ComplResult) {
	result, err := c.finder(matchResult, options.ReturnAll)
	if err != nil {
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

func (c *Completion) finder(match *MatchResult, returnAll bool) ([]string, error) {

	complResult, err := c.generateEntries(match, returnAll)
	if err != nil {
		return nil, err
	}
	if len(complResult) == 0 {
		return []string{}, nil
	}

	if len(complResult) == 1 || c.SelectOne {
		return []string{complResult[0]}, nil
	}
	if returnAll {
		return complResult, nil
	}
	return complResult, nil
}

func (c *Completion) renderLBuffer(matchR *MatchResult, items []string) (string, error) {
	if c.Replacement == "" {
		return matchR.line + strings.Join(items, c.ItemSeparator), nil
	} else if c.Command == "" || len(items) > 0 {
		vars := make(map[string]string)
		vars["item"] = strings.Join(items, c.ItemSeparator)
		for k, v := range matchR.subExpNamed {
			vars[k] = v
		}
		buf, err := render(c.Replacement, matchR.subExp, vars)
		if err != nil {
			return "", err
		}
		return matchR.line[:matchR.startIdx] + buf.String(), nil
	}

	return "", nil
}
