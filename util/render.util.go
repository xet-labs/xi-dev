package util

import (
	"bytes"
	"html/template"
)

func RenderPage(tmpl *template.Template, raw string, data any) (string, error) {
	// Clone the tmpl layout template to keep isolation between requests
	t, _ := tmpl.Clone()


	// Parse the content into the cloned template (defines "content")
	if _, err := t.Parse(raw); err != nil {
		return "", err
	}

	// Execute layout/tmpl using the merged data
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "layout/base", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}