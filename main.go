package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

var tmpl = template.Must(template.ParseGlob("view/page/*.*"))

func index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index", nil)
}

func rootpage(w http.ResponseWriter, r *http.Request) {
	receive := mux.Vars(r)
	page := receive["page"]
	match, err := filepath.Glob("./view/page/" + page + ".*")
	if err != nil {
		fmt.Println(err)
	}
	if match != nil {
		tmpl.ExecuteTemplate(w, page, nil)
	} else {
		fmt.Fprintf(w, "404 page not found")
	}
}

// blom
func firstpage(w http.ResponseWriter, r *http.Request) {
	receive := mux.Vars(r)
	first := receive["first"]
	page := receive["page"]
	if first != "" {
		match, err := filepath.Glob("./view/page/" + first + "/" + page + ".*")
		fmt.Println(match)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(match)
		if match != nil {
			var tmpls = template.Must(template.ParseGlob("view/page/" + first + "/*.*"))
			tmpls.ExecuteTemplate(w, "index", nil)
		} else {
			fmt.Fprintf(w, "404 page not found")
		}
	} else if first != "" && page != "" {
		match, err := filepath.Glob("./view/page/" + first + "/" + page + ".*")
		if err != nil {
			fmt.Println(err)
		}
		if match != nil {
			var tmpls = template.Must(template.ParseGlob("view/page/" + first + "/*.*"))
			tmpls.ExecuteTemplate(w, page, nil)
		} else {
			fmt.Fprintf(w, "404 page not found")
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/view/").Handler(http.StripPrefix("/view/", http.FileServer(http.Dir("view/"))))
	r.HandleFunc("/", index)
	r.HandleFunc("/{page}", rootpage)
	r.HandleFunc("/{first}/{page}", firstpage)
	fmt.Println("ready")
	http.ListenAndServe(":8000", r)
}
