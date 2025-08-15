package htmlfn

import (
	"html/template"
	"log"
	"strings"
)

func LinkLib(libs ...string) template.HTML {
	if len(libs) == 0 {
		return ""
	}

	var b strings.Builder

	for _, lib := range libs {
		switch lib {
		case "fa":
			b.WriteString(string(LinkCss("preload", "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css")))
		case "hljs":
			b.WriteString(string(LinkCss("preload", "/res/lib/hljs/styles/atom-one-dark.css")))
			b.WriteString(string(LinkJs("simple", "/res/lib/hljs/hl.min.js")))
			b.WriteString(`<script>hljs.highlightAll();</script>`)

		case "prism":
			b.WriteString(string(LinkCss("preload", "/res/lib/prism/prism.css")))
			b.WriteString(`<script defer src="/res/lib/prism/prism.js"></script>`)

		case "tw":
			b.WriteString(string(LinkCss("preload", "/res/lib/tailwind/tailwind.css")))

		case "mathjax":
			b.WriteString(`<script defer src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>`)

		default:
			log.Printf("[Htmlfn] WRN: Lib not found: %s", lib)
		}
	}

	return template.HTML(b.String())
}

func LinkLibSlice(libs []string) template.HTML { return LinkLib(libs...) }
