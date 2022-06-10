package compl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	assert := assert.New(t)
	m := matcher{pattern: "(foo) (?P<bar>bar) (xyz)"}
	var r *MatchResult
	var err error
	r, err = m.match("foo bar xyz")
	assert.NoError(err)
	assert.NotNil(r, "pattern should be matched")

	r, err = m.match("foo bar  xyz")
	assert.NoError(err)
	assert.Nil(r, "pattern does not match extra space")
}

func TestMatchReturnsRightMostResult(t *testing.T) {
	assert := assert.New(t)
	c := matcher{pattern: "([fb]oo)"}
	r, _ := c.match("foo boo")
	assert.Equal(4, r.startIdx, "starts at 'b'")
	assert.Equal(r.subExp[0], "boo")
}

func TestMatchNamedSubExps(t *testing.T) {
	assert := assert.New(t)
	c := matcher{pattern: "(?P<bar>bar)"}
	r, _ := c.match("bar xyz bar")
	assert.Equal(8, r.startIdx, "starts at the second 'b'")
	assert.Equal("bar", r.subExpNamed["bar"])
	assert.Equal(0, len(r.subExp))
}
