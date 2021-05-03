package main

import (
	"html/template"
	"log"
	"net/http"
)

var serverHandler *http.ServeMux
var staticHandler http.Handler
var server http.Server

var templates *template.Template

func main() {
	templates = template.Must(template.ParseFiles("index.html", "stats.html"))

	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":3000", Handler: serverHandler}

	staticHandler = http.FileServer(http.Dir("static"))
	serverHandler.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	serverHandler.HandleFunc("/", indexHandler)

	log.Println("server starting on Port :3000")
	log.Fatal(server.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}
