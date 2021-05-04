package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

//holds a single item when fetching the whole modList
type ModListItem struct {
	DisplayName        string
	Rank               int
	DownloadsTotal     int
	DownloadsToday     int
	DownloadsYesterday int
	TModLoaderVersion  string
	ModName            string
}

//holds mod info that is fetched from tmlapis.repl.co/modInfo?mod=
type ModInfo struct {
	DisplayName        string
	Rank               int //gets added seperatly
	InternalName       string
	Author             string
	Homepage           string
	Description        string
	Icon               string
	Version            string
	TModLoaderVersion  string
	LastUpdated        string
	ModDependencies    string
	ModSide            string
	DownloadLink       string
	DownloadsTotal     int
	DownloadsToday     int //gets added seperatly
	DownloadsYesterday int
}

var ModNameMap map[string]string = make(map[string]string)           //maps Display names to Internal names
var ModInfoMap map[string]ModListItem = make(map[string]ModListItem) //maps Internal names to ModInfo (for Rank and DownloadsToday)

var random *rand.Rand = rand.New(rand.NewSource(time.Now().Unix())) //random device

func updateModMaps() error {
	resp, err := http.Get("https://tmlapis.repl.co/modList")
	if err != nil {
		return err
	}

	var ModList []ModListItem
	err = json.NewDecoder(resp.Body).Decode(&ModList) //decode the modlist
	if err != nil {
		return err
	}

	for _, v := range ModList {
		ModNameMap[url.QueryEscape(v.DisplayName)] = v.ModName //map all Display names to Internal names
		ModInfoMap[v.ModName] = v                              //map all Internal names to ModInfo data
	}
	return nil
}

//helper to return a json encoded struct
func returnJsonFromStruct(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

//api handler to get the whole modList
func getModlistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method must be of type GET", http.StatusBadRequest)
		return
	}

	resp, err := http.Get("https://tmlapis.repl.co/modList") //fetch the data
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var ModList []ModListItem
	err = json.NewDecoder(resp.Body).Decode(&ModList) //decode the data
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	returnJsonFromStruct(w, ModList, http.StatusOK) //return it to the frontend
}

//returns the Internal name of the given Display name from a query
func getInternalNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method must be of type GET", http.StatusBadRequest)
		return
	}

	DisplayName := r.URL.Query().Get("displayname")                     //read the query
	name, err := json.Marshal(ModNameMap[url.QueryEscape(DisplayName)]) //find the Internal name in the map
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(name)
}

//api handler to get info to a single mod whichs Internal name is given in the query
func getModInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method must be of type GET", http.StatusBadRequest)
		return
	}
	modName := r.URL.Query().Get("modname") //get the query
	log.Println("recieved request for " + modName)

	resp, err := http.Get("https://tmlapis.repl.co/modInfo?modname=" + modName) //fetch most of the data
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var modInfo ModInfo
	err = json.NewDecoder(resp.Body).Decode(&modInfo) //encode the data (without rank and DownloadsToday)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	modInfo.Rank = ModInfoMap[modInfo.InternalName].Rank                     //add the rank from the map
	modInfo.DownloadsToday = ModInfoMap[modInfo.InternalName].DownloadsToday //add DownloadsToday from the map
	err = json.NewEncoder(w).Encode(modInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//returns a random Mod name
func getRandomModHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method must be GET", http.StatusBadRequest)
		return
	}
	count := random.Intn(len(ModNameMap))
	i := 0
	for _, v := range ModNameMap {
		if i >= count {
			name, err := json.Marshal(url.QueryEscape(v))
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(name)
			return
		}
		i++
	}
}
