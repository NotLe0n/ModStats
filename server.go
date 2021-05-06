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

var wg sync.WaitGroup //needed to keep the goroutines in sync when the server is shut down

var serverHandler *http.ServeMux //for all the request handler
var staticHandler http.Handler   //serves the static folder (javascript + css and assets)
var server http.Server           //the server itself

var templates *template.Template //the html files

//load the html files
func loadTemplates() {
	templates = template.Must(template.ParseFiles("index.html", "stats.html", "modList.html", "author.html"))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dataMutex = &sync.Mutex{}

	loadTemplates()

	//this goroutine updates the mod list every 10 minutes so that the loading time is not too long on every reload
	go func() {
		for ; true; <-time.Tick(10 * time.Minute) {
			log.Println("updating ModNameMap")
			err := updateModMaps()
			if err != nil {
				log.Fatal("Unable to update ModNameMap: " + err.Error())
			}
		}
	}()

	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":3000", Handler: serverHandler} //server on Port :3000

	staticHandler = http.FileServer(http.Dir("static")) //serves the static directory for js + css and assets
	serverHandler.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	//add the html handler
	serverHandler.HandleFunc("/", indexHandler)
	serverHandler.HandleFunc("/stats", statsHandler)
	serverHandler.HandleFunc("/modList", modListHandler)
	//add the api handler so the frontend can fetch all the data it needs
	serverHandler.HandleFunc("/api/getModlist", getModlistHandler)
	serverHandler.HandleFunc("/api/getModInfo", getModInfoHandler)
	serverHandler.HandleFunc("/api/getInternalName", getInternalNameHandler)
	serverHandler.HandleFunc("/api/getRandomMod", getRandomModHandler)
	serverHandler.HandleFunc("/api/getAuthorInfo", getAuthorInfoHandler)

	log.Println("Starting cmd goroutine")
	wg.Add(1)
	//this goroutine is for the cmd interface (at the moment only the quit command for a gracefull shutdown)
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

//commands in the switch can be run from the server console (not in sync with the logging)
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

	steam64ID := r.URL.Query().Get("author")
	if steam64ID != "" {
		log.Println("Someone searched for an author!")

		err := templates.ExecuteTemplate(w, "author.html", struct{ Steam64ID string }{
			Steam64ID: steam64ID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	modName := r.URL.Query().Get("mod")
	if modName != "" {
		log.Println("Someone searched for a mod!")

		err := templates.ExecuteTemplate(w, "stats.html", struct{ Mod string }{
			Mod: modName,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}
}

func modListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Someone visited the mod list!")
	loadTemplates()
	err := templates.ExecuteTemplate(w, "modList.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}
