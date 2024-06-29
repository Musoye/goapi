package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Proverb struct {
	Content   string
	CreatedOn time.Time
}

var proverbStore = make(map[string]Proverb)
var id int = 0

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("templates/index.html",
		"templates/base.html"))
	templates["add"] = template.Must(template.ParseFiles("templates/add.html",
		"templates/base.html"))
	templates["edit"] = template.Must(template.ParseFiles("templates/edit.html",
		"templates/base.html"))
}

func renderTemplate(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "The template does not exist.", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, template, viewModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveProverb(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	content := r.PostFormValue("content")
	proverb := Proverb{content, time.Now()}
	id++
	k := strconv.Itoa(id)
	proverbStore[k] = proverb
	http.Redirect(w, r, "/", 302)
}

func addProverb(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "add", "base", nil)
}

type EditProverb struct {
	Proverb
	Id string
}

func editProverb(w http.ResponseWriter, r *http.Request) {
	var viewModel EditProverb
	vars := mux.Vars(r)
	k := vars["id"]
	if proverb, ok := proverbStore[k]; ok {
		viewModel = EditProverb{proverb, k}
	} else {
		http.Error(w, "Could not find the resource to edit.", http.StatusBadRequest)
	}
	renderTemplate(w, "edit", "base", viewModel)
}

func updateProverb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	var proverbToUpd Proverb
	if proverb, ok := proverbStore[k]; ok {
		r.ParseForm()
		proverbToUpd.Content = r.PostFormValue("content")
		proverbToUpd.CreatedOn = proverb.CreatedOn
		delete(proverbStore, k)
		proverbStore[k] = proverbToUpd
	} else {
		http.Error(w, "Could not find the resource to update.", http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", 302)
}

func deleteProverb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	if _, ok := proverbStore[k]; ok {
		delete(proverbStore, k)
	} else {
		http.Error(w, "Could not find the resource to delete.", http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", 302)
}

func getProverbs(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", "base", proverbStore)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/", fs)
	r.HandleFunc("/", getProverbs)
	r.HandleFunc("/proverbs/add", addProverb)
	r.HandleFunc("/proverbs/save", saveProverb)
	r.HandleFunc("/proverbs/edit/{id}", editProverb)
	r.HandleFunc("/proverbs/update/{id}", updateProverb)
	r.HandleFunc("/proverbs/delete/{id}", deleteProverb)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
