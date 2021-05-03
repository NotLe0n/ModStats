package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type AuthorModInfo struct {
	DisplayName        string
	DownloadsTotal     int
	DownloadsYesterday int
	TModLoaderVersion  string
	ModName            string
}

type ModInfo struct {
	DisplayName        string
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
	DownloadsYesterday int
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
	var ModList []AuthorModInfo
	err = json.NewDecoder(resp.Body).Decode(&ModList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	returnJsonFromStruct(w, ModList, http.StatusOK)
}

func getModInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method must be of type GET", http.StatusBadRequest)
		return
	}
	modName := r.URL.Query().Get("modname")
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
	err = json.NewEncoder(w).Encode(modInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
