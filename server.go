package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

var serverHandler *http.ServeMux
var staticHandler http.Handler
var server http.Server

var templates *template.Template

func loadTemplates() {
	templates = template.Must(template.ParseFiles("index.html", "stats.html"))
}

func main() {
	loadTemplates()

	err := updateModNameMap()
	if err != nil {
		log.Fatal("Unable to update ModNameMap: " + err.Error())
	}

	go func() {
		for range time.Tick(10 * time.Minute) {
			log.Println("updating ModNameMap")
			err := updateModNameMap()
			if err != nil {
				log.Fatal("Unable to update ModNameMap: " + err.Error())
			}
		}
	}()

	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":3000", Handler: serverHandler}

	staticHandler = http.FileServer(http.Dir("static"))
	serverHandler.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	serverHandler.HandleFunc("/", indexHandler)
	serverHandler.HandleFunc("/stats", statsHandler)

	serverHandler.HandleFunc("/api/getModlist", getModlistHandler)
	serverHandler.HandleFunc("/api/getModInfo", getModInfoHandler)
	serverHandler.HandleFunc("/api/getInternalName", getInternalNameHandler)

	log.Println("server starting on Port :3000")
	log.Fatal(server.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	loadTemplates() //we reload the templates on each call, so that we don't need to restart the server when changing the html (mainly for debugging)
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	loadTemplates() //we reload the templates on each call, so that we don't need to restart the server when changing the html (mainly for debugging)
	modName := r.URL.Query().Get("mod")
	err := templates.ExecuteTemplate(w, "stats.html", struct{ Mod string }{
		Mod: modName,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}
