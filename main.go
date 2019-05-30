package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/oxtoacart/bpool"
)

func init() {

	// if templates == nil {
	// 	templates = make(map[string]*template.Template)
	// }
	// //templatesDir := "./templates/"
	// templatesDir := "./views/"

	// layouts, err := filepath.Glob(templatesDir + "layouts/admin.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// parts, err := filepath.Glob(templatesDir + "parts/*.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pages, err := filepath.Glob(templatesDir + "pages/*.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, page := range pages {
	// 	files := append(parts, page, layouts[0])
	// 	Base(page))
	// 	templates[filepath.Base(page)] = template.Must(template.ParseFiles(files...))
	// }

	// if err != nil {
	// 	log.Fatal(err)
	// }
	loadConfiguration()
	loadTemplates()
}

var templates map[string]*template.Template
var bufpool *bpool.BufferPool

type TemplateConfig struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
}

var mainTmpl = `{{define "main" }} {{ template "layout" . }} {{ end }}`

var templateConfig TemplateConfig

func loadConfiguration() {
	templateConfig.TemplateLayoutPath = "views/layouts/"
	templateConfig.TemplateIncludePath = "views/pages/"
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, err := filepath.Glob(templateConfig.TemplateLayoutPath + "*.html")
	if err != nil {
		log.Fatal(err)
	}

	includeFiles, err := filepath.Glob(templateConfig.TemplateIncludePath + "*.html")
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

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
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

func main() {

	server := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/", dealTpl)
	http.HandleFunc("/list", dealList)

	log.Fatal(server.ListenAndServe())

}

// func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
// 	tmpl, ok := templates[name]
// 	if !ok {
// 		return fmt.Errorf("The template %s does not exist.", name)
// 	}
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	return tmpl.ExecuteTemplate(w, name, data)
// }

type Page struct {
	Title string
	Body  string
	Datas []int
	Mul   interface{}
}

func mul(a, b int) int {
	return a * b
}

func dealList(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "list.html", "")

}

func dealTpl(w http.ResponseWriter, r *http.Request) {

	pdata := &Page{
		Title: "Deal tpl",
		Body:  "xxxx666",
		Datas: []int{1, 2, 3, 4, 5, 6, 8, 2, 1, 23, 3, 545, 1, 54, 818, 6, 5194, 1561, 6, 56},
		Mul:   mul,
	}

	renderTemplate(w, "index.html", pdata)

}
