package htmx

import (
	_ "embed"
	"text/template"
)

//go:embed base.html
var baseHTML string

const (
	baseTmplName = "base"
)

var baseTmpl = template.Must(
	template.Must(
		template.New("html").Funcs(TemplateFuncMap).
			New(baseTmplName).Parse(baseHTML),
	).New("nav").Parse(navHTML),
)
