package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type ModListItem struct {
	DisplayName        string
	Rank               int
	DownloadsTotal     int
	DownloadsToday     int
	DownloadsYesterday int
	TModLoaderVersion  string
	ModName            string
}

type ModInfo struct {
	DisplayName        string
	Rank               int
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
	DownloadsToday     int
	DownloadsYesterday int
}

var ModNameMap map[string]string = make(map[string]string)
var ModInfoMap map[string]ModListItem = make(map[string]ModListItem)

var random *rand.Rand = rand.New(rand.NewSource(time.Now().Unix()))

func updateModMaps() error {
	resp, err := http.Get("https://tmlapis.repl.co/modList")
	if err != nil {
		return err
	}

	var ModList []ModListItem
	err = json.NewDecoder(resp.Body).Decode(&ModList)
	if err != nil {
		return err
	}

	for _, v := range ModList {
		ModNameMap[url.QueryEscape(v.DisplayName)] = v.ModName
		ModInfoMap[v.ModName] = v
	}
	return nil
}

func returnJsonFromStruct(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func getModlistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method must be of type GET", http.StatusBadRequest)
		return
	}

	resp, err := http.Get("https://tmlapis.repl.co/modList")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var ModList []ModListItem
	err = json.NewDecoder(resp.Body).Decode(&ModList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	returnJsonFromStruct(w, ModList, http.StatusOK)
}

func getInternalNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method must be of type GET", http.StatusBadRequest)
		return
	}

	DisplayName := r.URL.Query().Get("displayname")
	name, err := json.Marshal(ModNameMap[url.QueryEscape(DisplayName)])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(name)
}

func getModInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method must be of type GET", http.StatusBadRequest)
		return
	}
	modName := r.URL.Query().Get("modname")
	log.Println("recieved request for " + modName)

	resp, err := http.Get("https://tmlapis.repl.co/modInfo?modname=" + modName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var modInfo ModInfo
	err = json.NewDecoder(resp.Body).Decode(&modInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	modInfo.Rank = ModInfoMap[modInfo.InternalName].Rank
	modInfo.DownloadsToday = ModInfoMap[modInfo.InternalName].DownloadsToday
	err = json.NewEncoder(w).Encode(modInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
