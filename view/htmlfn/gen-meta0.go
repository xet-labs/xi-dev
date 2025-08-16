package htmlfn

import (
	"fmt"
	"html"
	"strings"
	"xi/app/model"
)

// Helper to add a tag if value is present
func addTag(b *strings.Builder, value, format string) {
	if value != "" {
		fmt.Fprintf(b, format+"\n", html.EscapeString(value))
	}
}

func GenerateHTML(m model.PageMeta) string {
	var b strings.Builder

	// Title
	addTag(&b, m.Title, "<title>%s</title>")
	addTag(&b, m.Description, `<meta name="description" content="%s">`)
	addTag(&b, m.Canonical, `<link rel="canonical" href="%s">`)

	// Tags
	if len(m.Tags) > 0 {
		addTag(&b, strings.Join(m.Tags, ", "), `<meta name="keywords" content="%s">`)
	}

	addTag(&b, m.Robots, `<meta name="robots" content="%s">`)
	addTag(&b, m.Locale, `<meta property="og:locale" content="%s">`)

	// Alternate
	for _, alt := range m.Alternate {
		addTag(&b, alt, `<link rel="alternate" hreflang="%s">`)
	}

	// Open Graph
	addTag(&b, m.OG.Title, `<meta property="og:title" content="%s">`)
	addTag(&b, m.OG.Description, `<meta property="og:description" content="%s">`)
	addTag(&b, m.OG.Image, `<meta property="og:image" content="%s">`)
	addTag(&b, m.OG.Type, `<meta property="og:type" content="%s">`)
	addTag(&b, m.OG.URL, `<meta property="og:url" content="%s">`)

	// Twitter
	addTag(&b, m.Twitter.Card, `<meta name="twitter:card" content="%s">`)
	addTag(&b, m.Twitter.Title, `<meta name="twitter:title" content="%s">`)
	addTag(&b, m.Twitter.Description, `<meta name="twitter:description" content="%s">`)
	addTag(&b, m.Twitter.Image, `<meta name="twitter:image" content="%s">`)
	addTag(&b, m.Twitter.Site, `<meta name="twitter:site" content="%s">`)
	addTag(&b, m.Twitter.Creator, `<meta name="twitter:creator" content="%s">`)

	// Structured JSON-LD (no escaping)
	// if m.LD != "" {
	// 	fmt.Fprintf(&b, `<script type="application/ld+json">%s</script>`+"\n", m.Structured)
	// }

	return b.String()
}