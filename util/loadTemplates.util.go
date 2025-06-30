package util

import (
	"html/template"
	"log"
)

func LoadTemplates(dir string, ext string) *template.Template {
	files, err := GetFilesWithExt(dir, ext)
	if err != nil {
		log.Fatalf("template load error: %v", err)
	}

	tmpl := template.Must(template.New("").
		Funcs(HtmlFuncs).
		// .Funcs(strutil.Funcs).
		// .Funcs(timeutil.Funcs).
		ParseFiles(files...),
	)

	return tmpl
}
