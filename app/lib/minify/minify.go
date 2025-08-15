package minify

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	min3 "github.com/tdewolff/minify/v2"
	min3css "github.com/tdewolff/minify/v2/css"
	min3html "github.com/tdewolff/minify/v2/html"
)

type MinifyLib struct {
	Ecli    *gin.Engine
	Tcli    *template.Template
	RawTcli *template.Template
}

var (
	Minify   = &MinifyLib{}
	minifier = min3.New()

	// CSS Regex
	reComment     = regexp.MustCompile(`(?s)/\*.*?\*/`)
	reWhitespace  = regexp.MustCompile(`\s+`)
	reSpaceAround = regexp.MustCompile(`\s*([{}:;,])\s*`)
	reEmptyRule   = regexp.MustCompile(`[^{}]+\{\}`)
)

func init() {
	minifier.Add("text/css", &min3css.Minifier{
		// KeepCSS2: false,
	})

	minifier.Add("text/html", &min3html.Minifier{
		KeepWhitespace:          false, // collapse all whitespace
		KeepComments:            false, // remove comments
		KeepConditionalComments: false, // remove IE conditionals too
		KeepDefaultAttrVals:     false, // remove default values (e.g., type="text")
		KeepDocumentTags:        false, // remove <!DOCTYPE> if possible
		KeepEndTags:             false, // remove optional closing tags (</li>, etc.)
		KeepQuotes:              false, // remove attribute quotes where safe
	})
}

// Minify CSS
func (m *MinifyLib) Css(input []byte) ([]byte, error) {
	var out bytes.Buffer
	err := minifier.Minify("text/css", &out, bytes.NewReader(input))

	return out.Bytes(), err
}

// Minify CSS Hybrid
func (m *MinifyLib) CssHybrid(input []byte) ([]byte, error) {
	var out bytes.Buffer
	if err := minifier.Minify("text/css", &out, bytes.NewReader(input)); err != nil {
		return nil, err
	}

	// Remove empty rules recursively
	css := out.String()
	for {
		cleaned := reEmptyRule.ReplaceAllString(css, "")
		if cleaned == css {
			break
		}
		css = cleaned
	}
	return []byte(css), nil
}

// Minify CSS Regex
func (m *MinifyLib) CssRegex(css string) (string, error) {
	css = reComment.ReplaceAllString(css, "")       // Remove comments
	css = reWhitespace.ReplaceAllString(css, " ")   // Collapse spaces
	css = reSpaceAround.ReplaceAllString(css, "$1") // Minify spaces
	css = strings.ReplaceAll(css, ";}", "}")        // Remove trailing semicolons before }

	// Remove empty rules recursively
	for {
		newCSS := reEmptyRule.ReplaceAllString(css, "")
		if newCSS == css {
			break
		}
		css = newCSS
	}

	return strings.TrimSpace(css), nil
}

// Minify HTML
func (m *MinifyLib) Html(input []byte) ([]byte, error) {
	var out bytes.Buffer
	err := minifier.Minify("text/html", &out, bytes.NewReader(input))
	return out.Bytes(), err
}
