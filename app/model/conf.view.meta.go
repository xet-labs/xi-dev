package model

import "time"

type PageMeta struct {
	Type string `json:"type,omitempty"` // WebSite, BlogPosting, Product, NewsArticle, Article
	// Basic SEO
	Title       string   `json:"title,omitempty"`
	URL         string   `json:"url,omitempty"`
	Canonical   string   `json:"canonical,omitempty"`
	Description string   `json:"description,omitempty"`
	Tagline     string   `json:"tagline,omitempty"`
	Image      Image  `json:"images,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Hrefs       []HrefLang
	Locale      string   `json:"locale,omitempty"` // en_US etc (for og:locale)
	Robots      string `json:"robots,omitempty"` // e.g., "index, follow"
	Referrer    string // default "no-referrer-when-downgrade"

	// Social/owners
	Author    Author    `json:"author,omitempty"`
	Publisher Publisher `json:"publisher,omitempty"`
	OG        OG        `json:"og,omitempty"`
	Twitter   Twitter   `json:"twitter,omitempty"`

	// Article/Product specifics
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	ReadingTime string     `json:"readingTime,omitempty"` // e.g., "5 min"
	Category    string     `json:"category,omitempty"`
	IsFree      *bool      `json:"isAccessibleForFree,omitempty"`

	Alternate []string `json:"alternate,omitempty"` // hreflang alternate URLs

	// Extra JSON-LD (raw block to merge)
	LDPre  map[string]any `json:"ld_pre,omitempty"`
	LDPost map[string]any `json:"ld_post,omitempty"`
	LD     map[string]any `json:"ld,omitempty"`

	Brand PageBrand `json:"brand"`
}

type PageBrand struct {
	IncTitleSuffix *bool  `json:"inc_title_suffix,omitempty"`
	TitleSuffixSep string `json:"title_suffix_sep,omitempty"`
}

type Image struct {
	URL string `json:"url,omitempty"`
	Alt string `json:"alt,omitempty"`
}

type HrefLang struct {
	Lang string `json:"lang,omitempty"` // e.g. en, en-IN, fr
	URL  string `json:"url,omitempty"`
}

type Publisher struct {
	Name    string `json:"name,omitempty"`
	AltName string `json:"alt_name,omitempty"`
	URL     string `json:"url,omitempty"`
	Logo    string `json:"logo,omitempty"`
	LogoAlt string `json:"logo_alt,omitempty"`
}

type Author struct {
	Name        string `json:"name,omitempty"`
	URL         string `json:"url,omitempty"`
	Image       string `json:"image,omitempty"`
	JobTitle    string `json:"jobTitle,omitempty"`
	Description string `json:"description,omitempty"`
	SameAs      string `json:"sameAs,omitempty"` // single URL or CSV
}

type OG struct {
	Type        string            `json:"type,omitempty"`
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Image       string            `json:"image,omitempty"`
	URL         string            `json:"url,omitempty"`
	Extra       map[string]string `json:"extra,omitempty"`
}

type Twitter struct {
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Site        string            `json:"site,omitempty"`    // @handle
	Creator     string            `json:"creator,omitempty"` // @author
	Card        string            `json:"card,omitempty"`    // summary, summary_large_image
	Image       string            `json:"image,omitempty"`
	Extra       map[string]string `json:"extra,omitempty"` // label1/data1... or any kv
}
