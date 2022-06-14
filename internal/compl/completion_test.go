package compl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderLBuffer(t *testing.T) {
	const line = "xyz"

	c := Completion{ItemSeparator: " "}
	r, _ := c.renderLBuffer(&MatchResult{line: line}, []string{" bar", "baz"})
	assert.Equal(t, line+" bar baz", r)

}

func TestRenderSeparator(t *testing.T) {
	const line = "xyz"

	c := Completion{ItemSeparator: "\t"}
	r, _ := c.renderLBuffer(&MatchResult{line: line}, []string{" bar", "baz"})
	assert.Equal(t, line+" bar\tbaz", r)

}

func TestRenderLBufferWReplacement(t *testing.T) {
	const line = "before match xyz"

	c := Completion{ItemSeparator: " ", Replacement: "{{._1}} {{.item}} {{.foo}}"}
	r, _ := c.renderLBuffer(
		&MatchResult{line: line,
			startIdx:    13,
			subExp:      []string{"abc"},
			subExpNamed: map[string]string{"foo": "fooval"}},
		[]string{" bar", "baz"})
	assert.Equal(t, "before match abc  bar baz fooval", r)
}

func TestEmptyCompletion(t *testing.T) {
	c := NewCompletion("foo bar")
	r, code := c.MatchAndFind("foo bar", CompletionOptions{})
	assert.Equal(t, "foo bar", r[0])
	assert.Equal(t, Matched, code)
}

func TestEmptyCompletionWReplacement(t *testing.T) {
	c := NewCompletion("(foo) bar")
	c.Replacement = "{{ ._1 }} xyz"
	r, code := c.MatchAndFind("foo bar", CompletionOptions{})
	assert.Equal(t, "foo xyz", r[0])
	assert.Equal(t, Matched, code)
}

func TestCompletionShellTmpl(t *testing.T) {
	c := NewCompletion("(foo) bar")
	c.Replacement = `{{ shell "printf %s" ._1 }} xyz`
	r, code := c.MatchAndFind("foo bar", CompletionOptions{})
	assert.Equal(t, "foo xyz", r[0])
	assert.Equal(t, Matched, code)
}

func TestCompletionCmdTmpl(t *testing.T) {
	c := NewCompletion("(foo) bar")
	c.Replacement = `{{ cmd "printf" "%s" ._1 }} xyz`
	r, code := c.MatchAndFind("foo bar", CompletionOptions{})
	assert.Equal(t, "foo xyz", r[0])
	assert.Equal(t, Matched, code)

	c.Replacement = `{{  ._1 | cmdPipe "tr" "o" "a" }} xyz`
	r, code = c.MatchAndFind("foo bar", CompletionOptions{})
	assert.Equal(t, "faa xyz", r[0])
	assert.Equal(t, Matched, code)
}

func TestCompletionReturnAll(t *testing.T) {
	c := NewCompletion("foo bar")
	c.Command = "echo \"xyz\nxyz\nxyz\""
	r, code := c.MatchAndFind("foo bar", CompletionOptions{ReturnAll: true})
	assert.Equal(t, 3, len(r))
	assert.Equal(t, Matched, code)
}

func TestCompletionUseArgsInCmd(t *testing.T) {
	c := NewCompletion("(foo) (?P<baz>bar)")
	c.Command = "echo $1 $baz"
	r, code := c.MatchAndFind("foo bar", CompletionOptions{ReturnAll: true})
	assert.Equal(t, 1, len(r))
	assert.Equal(t, "foo bar", r[0])
	assert.Equal(t, Matched, code)
}
