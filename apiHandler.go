package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ModInfo struct {
	DisplayName        string
	DownloadsTotal     int
	DownloadsYesterday int
	TModLoaderVersion  string
	ModName            string
}

func returnJsonFromStruct(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func getModlistHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://tmlapis.repl.co/modList")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	var ModList []ModInfo
	err = json.NewDecoder(resp.Body).Decode(&ModList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	returnJsonFromStruct(w, ModList, http.StatusOK)
}
