package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

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
	serverHandler.HandleFunc("/api/getRandomMod", getRandomModHandler)

	log.Println("Starting cmd goroutine")
	wg.Add(1)
	go func() {
		defer wg.Done() //tell the waiter group that we are finished at the end
		cmdInterface()
		log.Println("cmd goroutine finished")
	}()

	log.Println("server starting on Port :3000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	} else if err == http.ErrServerClosed {
		log.Println("Server not listening anymore")
	}
}

func cmdInterface() {
	for loop := true; loop; {
		var inp string
		_, err := fmt.Scanln(&inp)
		if err != nil {
			log.Println(err.Error())
		} else {
			switch inp {
			case "quit":
				log.Println("Attempting to shutdown server")
				err := server.Shutdown(context.Background())
				if err != nil {
					log.Fatal("Error while trying to shutdown server: " + err.Error())
				}
				log.Println("Server was shutdown")
				loop = false
			default:
				fmt.Println("cmd not supported")
			}
		}
	}
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
