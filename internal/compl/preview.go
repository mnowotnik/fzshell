package compl

import (
	"text/template"

	"github.com/pkg/errors"
)

type Preview struct {
	PreviewTmpl *template.Template
	MatchResult *MatchResult
}

func (p *CompletionEntry) RenderPreview() (string, error) {
	if p.Preview == nil {
		return "", nil
	}
	tmplData := make(map[string]string)
	for k, v := range p.Preview.MatchResult.subExpNamed {
		tmplData[k] = v
	}
	tmplData["item"] = p.Item
	if buf, err := renderFromTemplate(p.Preview.PreviewTmpl, p.Preview.MatchResult.subExp, tmplData); err != nil {
		return "", errors.Wrapf(err, "error rendering template [item=%s]: %s", p.Item)
	} else {
		return buf.String(), nil
	}
}
