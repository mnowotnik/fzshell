package compl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_render(t *testing.T) {
	got, _ := render("foo {{.a}} {{._1}} {{.b}} xyz", []string{"1"}, map[string]string{"a": "a", "b": "b"})
	assert.Equal(t, "foo a 1 b xyz", got.String())
}
