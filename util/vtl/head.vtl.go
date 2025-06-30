package vtl

import (
	"html/template"
	"strings"
)

func Csslink(url string, load ...int) template.HTML {
	ld := 0
	if len(load) > 0 {
		ld = load[0]
	}

	escapedURL := template.HTMLEscapeString(url)
	var b strings.Builder

	switch ld {
	case 0:
		b.WriteString(`<link rel="preload" as="style" href="`)
		b.WriteString(escapedURL)
		b.WriteString(`"><link rel="stylesheet" href="`)
		b.WriteString(escapedURL)
		b.WriteString(`">`)
	case 1:
		b.WriteString(`<link rel="stylesheet" href="`)
		b.WriteString(escapedURL)
		b.WriteString(`">`)
	}

	return template.HTML(b.String())
}

func Jslink(url string, load ...int) template.HTML {
	ld := 0
	if len(load) > 0 {
		ld = load[0]
	}

	escapedURL := template.HTMLEscapeString(url)
	var b strings.Builder

	switch ld {
	case 0:
		b.WriteString(`<link rel="preload" as="script" href="`)
		b.WriteString(escapedURL)
		b.WriteString(`" crossorigin="anonymous">`)
		b.WriteString(`<script defer src="`)
		b.WriteString(escapedURL)
		b.WriteString(`" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`)
	case 1:
		b.WriteString(`<script src="`)
		b.WriteString(escapedURL)
		b.WriteString(`"></script>`)
	}

	return template.HTML(b.String())
}
