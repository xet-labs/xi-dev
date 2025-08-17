package view

import (
	"maps"
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"

	"xi/app/lib/cfg"
	"xi/app/model"
)

// <link rel="alternate" type="application/json+oembed" href="/oembed?url=https://mysite.com/page" />
// <link rel="alternate" type="application/json" href="/api/posts/123.json" />
// <link rel='shortlink' href='https://blueminch.com/' />

func (v *ViewLib) GenMeta(m model.PageMeta) string {

	var b strings.Builder

	// computed
	title := metaTitle(m)
	metaType(&m)

	// --- Basic meta ---
	add(&b, `<title>%s</title>`, title)
	add(&b, `<meta name="robots" content="%s">`, m.Robots)
	add(&b, `<meta name="referrer" content="%s">`, m.Referrer)
	add(&b, `<link rel="canonical" href="%s">`, m.URL)
	add(&b, `<meta name="title" content="%s">`, title)
	add(&b, `<meta name="description" content="%s">`, m.Description)
	if len(m.Tags) > 0 {
		add(&b, `<meta name="keywords" content="%s">`, strings.Join(m.Tags, ", "))
	}

	// hreflang alternates
	// for _, h := range m.Hrefs {
	// 	if h.Lang != "" && h.URL != "" {
	// 		fmt.Fprintf(&b, `<link rel="alternate" hreflang="%s" href="%s">\n`, e(h.Lang), e(h.URL))
	// 	}
	// }

	// --- Open Graph ---
	add(&b, `<meta property="og:type" content="%s">`, m.OG.Type)
	add(&b, `<meta property="og:locale" content="%s">`, m.Locale)
	add(&b, `<meta property="og:site_name" content="%s">`, cfg.Brand.Name)
	add(&b, `<meta property="og:title" content="%s">`, StrFallback(m.OG.Title, title))
	add(&b, `<meta property="og:description" content="%s">`, StrFallback(m.OG.Description, m.Description))
	add(&b, `<meta property="og:url" content="%s">`, StrFallback(m.OG.URL, m.URL))
	add(&b, `<meta property="og:image" content="%s">`, m.Img.URL)
	add(&b, `<meta property="og:image:alt" content="%s">`, StrFallback(m.Img.Alt, m.Title))

	// --- Twitter ---
	add(&b, `<meta name="twitter:card" content="%s">`, StrNotEmptyThen("summary_large_image", m.Img.URL))
	add(&b, `<meta name="twitter:site" content="%s">`, m.Twitter.Site)
	add(&b, `<meta name="twitter:creator" content="%s">`, m.Twitter.Creator)
	add(&b, `<meta property="twitter:domain" value="%s">`, cfg.Brand.Domain)
	add(&b, `<meta name="twitter:title" content="%s">`, StrFallback(m.OG.Title, title))
	add(&b, `<meta name="twitter:description" content="%s">`, StrFallback(m.OG.Description, m.Description))
	add(&b, `<meta name="twitter:url" content="%s">`, StrFallback(m.OG.URL, m.URL))
	add(&b, `<meta name="twitter:image:src" content="%s">`, m.Img.URL)
	add(&b, `<meta name="twitter:image:alt" content="%s">`, StrFallback(m.Img.Alt, m.Title))

	// Twitter extra labels (label1/data1 ... or arbitrary kv)
	if len(m.Twitter.Extra) > 0 {
		// If user provided label/data pairs, just emit them in order of key
		for k, v := range m.Twitter.Extra {
			fmt.Fprintf(&b, `<meta name="%s" content="%s">`+"\n", e(k), e(v))
		}
	}
	if isArticle(m.Type) {
		add(&b, `<meta name="twitter:label1" content="%s">`, StrNotEmptyThen("Written by", m.Author.Name))
		add(&b, `<meta name="twitter:data1" content="%s">`, "@" + m.Author.Name)

		add(&b, `<meta name="twitter:label2" content="%s">`, StrNotEmptyThen("Category", m.Category))
		add(&b, `<meta name="twitter:data2" content="%s">`, m.Category)

		add(&b, `<meta name="twitter:label3" content="%s">`, PtrNotNilThen("Published on", m.CreatedAt))
		add(&b, `<meta name="twitter:data3" content="%s">`, m.CreatedAt.UTC().Format(time.RFC3339Nano))

		add(&b, `<meta name="twitter:label4" content="%s">`, StrNotEmptyThen("Reading time", m.ReadingTime))
		add(&b, `<meta name="twitter:data4" content="%s">`, m.ReadingTime)
	}

	// --- Author rel links (article-like) ---
	if isArticle(m.Type) {
		add(&b, `<meta name="author" content="%s">`, m.Author.Name)
		add(&b, `<meta property="article:author" content="%s">`, m.Author.URL)
		add(&b, `<link rel="author" href="%s">`, m.Author.URL)
		if m.CreatedAt != nil {
			add(&b, `<meta property="article:published_time" content="%s">`, m.CreatedAt.UTC().Format(time.RFC3339Nano))
		}
	}

	// --- JSON-LD ---
	ld := ldJSON(m)
	if len(ld) > 0 {
		b.WriteString(`<script type="application/ld+json">`)
		b.Write(ld)
		b.WriteString(`</script>`)
	}

	return b.String()
}

// ------- JSON-LD builder -------

func ldJSON(m model.PageMeta) []byte {
	// computed
	title := metaTitle(m)

	root := map[string]any{
		"@context":         "https://schema.org",
		"@type":            m.Type,
		"url":              m.URL,
		"mainEntityOfPage": m.URL,
	}

	addIfNotEmptyTo(root, "description", m.Description)
	addIfNotEmptyTo(root, "image", m.Author.Image)
	addSliceIfNotEmpty(root, "keyword", m.Tags)

	// Article-like specifics
	if isArticle(m.Type) {
		addIfNotEmptyTo(root, "headline", title)
		addIfNotEmptyTo(root, "dateCreated", m.CreatedAt.UTC().Format(time.RFC3339Nano))
		addIfNotEmptyTo(root, "datePublished", m.CreatedAt.UTC().Format(time.RFC3339Nano))
		addIfNotEmptyTo(root, "dateModified", m.UpdatedAt.UTC().Format(time.RFC3339Nano))
		addSliceIfNotEmpty(root, "articleSection", m.Tags)
		
		root["author"] = map[string]any{} // create new map
    	authorMap := root["author"].(map[string]any)
		authorMap["@type"] = "Person"
		addIfNotEmptyTo(authorMap, "name", m.Author.Name)
		addIfNotEmptyTo(authorMap, "description", m.Author.Description)
		addIfNotEmptyTo(authorMap, "image", m.Author.Image)
		addIfNotEmptyTo(authorMap, "jobTitle", m.Author.JobTitle)
		addIfNotEmptyTo(authorMap, "sameAs", m.Author.SameAs)
		addIfNotEmptyTo(authorMap, "url", m.Author.URL)

		// convenience mirrors
		addIfNotEmptyTo(root, "creator", m.Author.Name)
		addIfNotEmptyTo(root, "editor", m.Author.Name)
	}

	// Access
	if m.IsFree != nil {
		root["isAccessibleForFree"] = *m.IsFree
	} else {
		root["isAccessibleForFree"] = true
	}

	// Publisher
	pub := map[string]any{}
	addIfNotEmptyTo(pub, "name", m.Publisher.Name)
	addIfNotEmptyTo(pub, "url", m.Publisher.URL)
	addIfNotEmptyTo(pub, "alternateName", m.Publisher.AltName)

	if m.Publisher.Logo != "" {
		pub["logo"] = map[string]any{"@type": "ImageObject", "url": m.Publisher.Logo}
	}
	if len(pub) > 0 {
		pub["@type"] = "Organization"
		root["publisher"] = pub
	}

	// Allow raw LD to merge (shallow)
	maps.Copy(root, m.LD)

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(root)
	return buf.Bytes()
}

// metaTitle builds the final page <title>
func metaTitle(m model.PageMeta) string {
	base := m.Title

	// add author if provided
	if isArticle(m.Type) && m.Author.Name != "" {
		base = base + " | by " + m.Author.Name
	}

	// append brand suffix if enabled
	if m.Brand.IncTitleSuffix != nil && *m.Brand.IncTitleSuffix {
		if m.Publisher.Name != "" {
			if m.Brand.TitleSuffixSep == "" {
				m.Brand.TitleSuffixSep = " | "
			}
			return base + m.Brand.TitleSuffixSep + m.Publisher.Name
		}
	}

	return base
}

func metaType(m *model.PageMeta) {
	if m.Type == "" {
		m.Type = "WebSite"
	}
	switch strings.ToLower(StrFallback(m.OG.Type, m.Type)) {
	case "article", "blogposting", "newsarticle":
		m.OG.Type = "article"
	case "product":
		m.OG.Type ="product"
	case "profile":
		m.OG.Type ="profile"
	default:
		m.OG.Type = "website"
	}
}

func isArticle(t string) bool {
	switch strings.ToLower(t) {
	case "blogposting", "newsarticle", "article":
		return true
	}
	return false
}

func boolOr(a, def bool) bool {
	if a {
		return true
	}
	return def
}

func addIfNotEmptyTo(m map[string]any, key string, val string) {
	if val != "" {
		m[key] = val
	}
}

func addSliceIfNotEmpty(m map[string]any, key string, vals []string) {
	if len(vals) > 0 {
		m[key] = vals
	}
}

func PtrNotNilThen[T any](val string, ptr *T) string {
	if ptr != nil {
		return val
	}
	return ""
}
func StrNotEmptyThen(val, str string) string {
	if str != "" {
		return val
	}
	return ""
}

func StrFallback(vals ...string) string {
	for _, v := range vals {
		if v == "" {
			continue
		}
		return v
	}
	return ""
}

func add(b *strings.Builder, f, v string) {
	if v != "" {
		fmt.Fprintf(b, f+"\n", e(v))
	}
}

func e(s string) string { return html.EscapeString(s) }
