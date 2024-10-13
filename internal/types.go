package internal

import "html/template"

type PageData struct {
	Body        template.HTML
	ProjectPath string
}
