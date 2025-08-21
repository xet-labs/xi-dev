package htmlfn

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"maps"
	"strings"
	"time"

	"xi/app/lib/util"
	"xi/app/lib/cfg"
	"xi/app/model"
)

func GenMeta(m *model.PageMeta) template.HTML {

	var b strings.Builder

	// computed
	title := metaTitle(m)
	metaType(m)

	// --- Basic meta ---
	add(&b, `<title>%s</title>`, title)
	add(&b, `<link rel="canonical" href="%s">`, m.URL)
	add(&b, `<link rel="shortlink" href="%s">`, m.ShortLink)
	add(&b, `<link rel="alternate" type="application/json" href="%s">`, m.AltJson)
	add(&b, `<meta name="title" content="%s">`, title)
	add(&b, `<meta name="description" content="%s">`, m.Description)
	add(&b, `<meta name="robots" content="%s">`, m.Robots)
	add(&b, `<meta name="referrer" content="%s">`, m.Referrer)
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
	add(&b, `<meta property="og:title" content="%s">`, util.Str.Fallback(m.OG.Title, title))
	add(&b, `<meta property="og:description" content="%s">`, util.Str.Fallback(m.OG.Description, m.Description))
	add(&b, `<meta property="og:url" content="%s">`, util.Str.Fallback(m.OG.URL, m.URL))
	add(&b, `<meta property="og:image" content="%s">`, m.Img.URL)
	add(&b, `<meta property="og:image:alt" content="%s">`, util.Str.Fallback(m.Img.Alt, m.Title))

	// --- Twitter ---
	add(&b, `<meta name="twitter:card" content="%s">`, util.Str.NotEmptyThen("summary_large_image", m.Img.URL))
	add(&b, `<meta name="twitter:site" content="%s">`, m.Twitter.Site)
	add(&b, `<meta name="twitter:creator" content="%s">`, m.Twitter.Creator)
	add(&b, `<meta property="twitter:domain" value="%s">`, cfg.Brand.Domain)
	add(&b, `<meta name="twitter:title" content="%s">`, util.Str.Fallback(m.OG.Title, title))
	add(&b, `<meta name="twitter:description" content="%s">`, util.Str.Fallback(m.OG.Description, m.Description))
	add(&b, `<meta name="twitter:url" content="%s">`, util.Str.Fallback(m.OG.URL, m.URL))
	add(&b, `<meta name="twitter:image:src" content="%s">`, m.Img.URL)
	add(&b, `<meta name="twitter:image:alt" content="%s">`, util.Str.Fallback(m.Img.Alt, m.Title))

	// Twitter extra labels (label1/data1 ... or arbitrary kv)
	if len(m.Twitter.Extra) > 0 {
		// If user provided label/data pairs, just emit them in order of key
		for k, v := range m.Twitter.Extra {
			fmt.Fprintf(&b, `<meta name="%s" content="%s">`+"\n", e(k), e(v))
		}
	}
	if isArticle(m.Type) {
		add(&b, `<meta name="twitter:label1" content="%s">`, util.Str.NotEmptyThen("Written by", m.Author.Name))
		add(&b, `<meta name="twitter:data1" content="%s">`, "@"+m.Author.Name)

		add(&b, `<meta name="twitter:label2" content="%s">`, util.Str.NotEmptyThen("Category", m.Category))
		add(&b, `<meta name="twitter:data2" content="%s">`, m.Category)

		add(&b, `<meta name="twitter:label3" content="%s">`, util.StrIfPtrNotNil("Published on", &m.CreatedAt))
		add(&b, `<meta name="twitter:data3" content="%s">`, m.CreatedAt.UTC().Format(time.RFC3339Nano))

		add(&b, `<meta name="twitter:label4" content="%s">`, util.Str.NotEmptyThen("Reading time", m.ReadingTime))
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

	if ld := ldJSON(m); len(ld) > 0 {
		b.WriteString(`<script type="application/ld+json">`)
		b.Write(ld)
		b.WriteString(`</script>`)
	}

	return template.HTML(b.String())
}

// ------- JSON-LD builder -------

func ldJSON(m *model.PageMeta) []byte {
	// computed
	title := metaTitle(m)

	root := map[string]any{
		"@context":         "https://schema.org",
		"@type":            m.Type,
		"url":              m.URL,
		"mainEntityOfPage": m.URL,
	}

	util.Map.AddIfNotEmpty(root, "description", m.Description)
	util.Map.AddIfNotEmpty(root, "image", m.Author.Img)
	util.Map.AddIfNotEmptySlice(root, "keywords", m.Tags)

	// Article-like specifics
	if isArticle(m.Type) {
		util.Map.AddIfNotEmpty(root, "headline", title)
		util.Map.AddIfNotEmpty(root, "dateCreated", m.CreatedAt.UTC().Format(time.RFC3339Nano))
		util.Map.AddIfNotEmpty(root, "datePublished", m.CreatedAt.UTC().Format(time.RFC3339Nano))
		util.Map.AddIfNotEmpty(root, "dateModified", m.UpdatedAt.UTC().Format(time.RFC3339Nano))
		util.Map.AddIfNotEmptySlice(root, "articleSection", m.Tags)

		root["author"] = map[string]any{} // create new map
		authorMap := root["author"].(map[string]any)
		authorMap["@type"] = "Person"
		util.Map.AddIfNotEmpty(authorMap, "name", m.Author.Name)
		util.Map.AddIfNotEmpty(authorMap, "description", m.Author.Description)
		util.Map.AddIfNotEmpty(authorMap, "image", m.Author.Img)
		util.Map.AddIfNotEmpty(authorMap, "jobTitle", m.Author.JobTitle)
		util.Map.AddIfNotEmpty(authorMap, "sameAs", m.Author.SameAs)
		util.Map.AddIfNotEmpty(authorMap, "url", m.Author.URL)

		// convenience mirrors
		util.Map.AddIfNotEmpty(root, "creator", m.Author.Name)
		util.Map.AddIfNotEmpty(root, "editor", m.Author.Name)
	}

	// Access
	if m.IsFree != nil {
		root["isAccessibleForFree"] = *m.IsFree
	} else {
		root["isAccessibleForFree"] = true
	}

	// Publisher
	pub := map[string]any{}
	util.Map.AddIfNotEmpty(pub, "name", m.Publisher.Name)
	util.Map.AddIfNotEmpty(pub, "url", m.Publisher.URL)
	util.Map.AddIfNotEmpty(pub, "alternateName", m.Publisher.AltName)

	if m.Publisher.Logo != "" {
		pub["logo"] = map[string]any{"@type": "ImageObject", "url": m.Publisher.Logo}
	}
	if len(pub) > 0 {
		pub["@type"] = "Organization"
		root["publisher"] = pub
	}

	// Allow raw LD to merge (shallow)
	maps.Copy(root, m.LD)

	if b, err := json.Marshal(root); err == nil {
		return b
	}
	return []byte{}
}

// metaTitle builds the final page <title>
func metaTitle(m *model.PageMeta) string {
	base := m.Title

	// add author if provided
	if isArticle(m.Type) && m.Author.Name != "" {
		base = base + " | by " + m.Author.Name
	}

	// append brand suffix if enabled
	if m.Brand.IncTitleSuffix != nil && *m.Brand.IncTitleSuffix {
		if m.Publisher.Name != "" {
			return base + m.Brand.TitleSuffixSep + m.Publisher.Name
		}
	}

	return base
}

func metaType(m *model.PageMeta) {
	if m.Type == "" {
		m.Type = "WebSite"
	}
	switch strings.ToLower(util.Str.Fallback(m.OG.Type, m.Type)) {
	case "article", "blogposting", "newsarticle":
		m.OG.Type = "article"
	case "product":
		m.OG.Type = "product"
	case "profile":
		m.OG.Type = "profile"
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

func add(b *strings.Builder, f, v string) {
	if v != "" {
		fmt.Fprintf(b, f+"\n", e(v))
	}
}

func e(s string) string { return html.EscapeString(s) }
