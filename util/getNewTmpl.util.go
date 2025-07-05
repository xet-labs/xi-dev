package util

import (
	"html/template"
	"log"
)

func GetNewTmpl(tmplName, ext string, dirs ...string) *template.Template {
	files, err := GetFilesWithExt(ext, dirs...)
	if err != nil {
		log.Fatalf("Route: template load error: %v", err)
	}

	tmpl := template.Must(template.New(tmplName).
		Funcs(HtmlFuncs).
		// .Funcs(strutil.Funcs).
		// .Funcs(timeutil.Funcs).
		ParseFiles(files...),
	)

	return tmpl
}