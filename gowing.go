package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/oxtoacart/bpool"
)

//Gowing struct
type Gowing struct {
}

var templates map[string]*template.Template
var bufpool *bpool.BufferPool

//TemplateConfig config data
type TemplateConfig struct {
	TemplateLayoutPath string
	TemplatePagesPath  string
	TemplateExtension  string
}

var mainTmpl = `{{define "main" }} {{ template "layout" . }} {{ end }}`

var tplConfig TemplateConfig

// func main() {
// 	var gow *Gowing
// 	gow.New("views/layouts/", "views/pages/", "html")

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		gow.View(w, "list.html", "")
// 	})

// 	http.ListenAndServe(":8080", nil)

// }

func main() {

}

//New a gowing
func (gw *Gowing) New(layoutPath, pagePath, extension string) {
	tplConfig.TemplateLayoutPath = layoutPath
	tplConfig.TemplatePagesPath = pagePath
	tplConfig.TemplateExtension = "*." + extension
	gw.LoadTemplates()

}

//LoadConfiguration load template config like view/layout/ view/pages/  html
func (gw *Gowing) LoadConfiguration(layoutPath, pagePath, extension string) {
	tplConfig.TemplateLayoutPath = layoutPath
	tplConfig.TemplatePagesPath = pagePath
	tplConfig.TemplateExtension = "*." + extension
}

//LoadTemplates load pages
func (gw *Gowing) LoadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, err := filepath.Glob(tplConfig.TemplateLayoutPath + tplConfig.TemplateExtension)
	if err != nil {
		log.Fatal(err)
	}

	includeFiles, err := filepath.Glob(tplConfig.TemplatePagesPath + tplConfig.TemplateExtension)
	if err != nil {
		log.Fatal(err)
	}

	mainTemplate := template.New("main")

	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range includeFiles {
		fileName := filepath.Base(file)

		files := append(layoutFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	log.Println("templates loading successful")

	bufpool = bpool.NewBufferPool(64)

	log.Println("buffer allocation successful")
}

//View show a html
func (gw *Gowing) View(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", name),
			http.StatusInternalServerError)
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.Execute(buf, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
