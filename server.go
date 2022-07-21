package main

import (
	"html/template"
	"log"
	"net/http"
)

var serverHandler *http.ServeMux //for all the request handler
var staticHandler http.Handler   //serves the static folder (javascript + css and assets)
var server http.Server           //the server itself

var templates *template.Template //the html files

var errLog *log.Logger

//load the html files
func loadTemplates() {
	templates = template.Must(template.ParseFiles("index.html", "stats.html", "modList.html", "author.html"))
}

/*func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	errLog = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	errLog.Println("Setup log for this Session")

	dataMutex = &sync.Mutex{}

	loadTemplates()

	//this goroutine updates the mod list every 10 minutes so that the loading time is not too long on every reload
	go func() {
		for ; true; <-time.Tick(10 * time.Minute) {
			log.Println("updating ModNameMap")
			err := updateModMaps()
			if err != nil {
				log.Println("Unable to update ModNameMap, using the last valid state: " + err.Error())
				errLog.Println("Unable to update ModNameMap: " + err.Error())
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
	serverHandler.HandleFunc("/api/getVersionHistory", getVersionHistory)
	serverHandler.HandleFunc("/api/getRandomMod", getRandomModHandler)
	serverHandler.HandleFunc("/api/getAuthorInfo", getAuthorInfoHandler)

	log.Println("Starting cmd goroutine")

	log.Println("server starting on Port :3000")
	/*if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	} else if err == http.ErrServerClosed {
		log.Println("Server not listening anymore")
	}
}*/

func statsHandler(w http.ResponseWriter, r *http.Request) {
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
