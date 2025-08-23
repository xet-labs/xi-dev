package htmlfn

import (
	"html/template"
	"strings"
)

func LinkCss(mode string, urls ...string) template.HTML {
	var b strings.Builder

	for _, u := range urls {
		escapedURL := template.HTMLEscapeString(u)

		switch mode {
		case "preload":
			b.WriteString(`<link rel="preload" as="style" href="`)
			b.WriteString(escapedURL)
			b.WriteString(`">`)
			b.WriteString(`<link rel="stylesheet" href="`)
			b.WriteString(escapedURL)
			b.WriteString(`">`)
		case "simple":
			fallthrough
		default:
			b.WriteString(`<link rel="stylesheet" href="`)
			b.WriteString(escapedURL)
			b.WriteString(`">`)
		}
	}

	return template.HTML(b.String())
}

func LinkCssSlice(mode string, slice []string) template.HTML { return LinkCss(mode, slice...) }

func LinkJs(mode string, urls ...string) template.HTML {
	var b strings.Builder

	for _, u := range urls {
		escapedURL := template.HTMLEscapeString(u)

		switch mode {
		case "preload":
			b.WriteString(`<link rel="preload" as="script" href="`)
			b.WriteString(escapedURL)
			b.WriteString(`" crossorigin="anonymous">`)
			b.WriteString(`<script defer src="`)
			b.WriteString(escapedURL)
			b.WriteString(`" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`)
		case "simple":
			fallthrough
		default:
			b.WriteString(`<script src="`)
			b.WriteString(escapedURL)
			b.WriteString(`"></script>`)
		}
	}

	return template.HTML(b.String())
}

func LinkJsSlice(mode string, slice []string) template.HTML { return LinkJs(mode, slice...) }
