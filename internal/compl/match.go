package compl

import (
	"regexp"
)

type MatchResult struct {
	subExp      []string
	subExpNamed map[string]string
	startIdx    int
	line        string
}

type matcher struct {
	matchMultipleSpaces bool
	pattern             string
}

func newMatchResult(startIdx int, line string) *MatchResult {
	return &MatchResult{
		subExp:      []string{},
		subExpNamed: make(map[string]string),
		line:        line,
		startIdx:    startIdx,
	}
}

func (m matcher) match(line string) (*MatchResult, error) {
	var pattern *regexp.Regexp
	var err error

	if m.matchMultipleSpaces {
		pattern, err = regexp.Compile(spaces.ReplaceAllString(m.pattern, " +") + "$")
	} else {
		pattern, err = regexp.Compile(m.pattern + "$")
	}
	if err != nil {
		return nil, err
	}
	subMatches := pattern.FindStringSubmatch(line)
	if len(subMatches) == 0 {
		return nil, nil
	}
	subExpMap := map[int]string{}
	for _, name := range pattern.SubexpNames()[1:] {
		subExpMap[pattern.SubexpIndex(name)] = name
	}
	result := newMatchResult(len(line)-len(subMatches[0]), line)
	for i, match := range subMatches[1:] {
		if _, ok := subExpMap[i+1]; ok {
			result.subExpNamed[subExpMap[i+1]] = match
		} else {
			result.subExp = append(result.subExp, match)
		}
	}
	return result, nil
}
