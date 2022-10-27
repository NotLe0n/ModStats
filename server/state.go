package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"sync"
)

// ModListItem holds a single item when fetching the whole modList
type ModListItem struct {
	DisplayName        string `json:"display_name"`
	Rank               int    `json:"rank"`
	DownloadsTotal     int    `json:"downloads_total"`
	DownloadsToday     int    `json:"downloads_today"`
	DownloadsYesterday int    `json:"downloads_yesterday"`
	TModLoaderVersion  string `json:"tmodloader_version"`
	InternalName       string `json:"internal_name"`
}

// ModInfo holds mod info that is fetched from tmlapis.repl.co/modInfo?mod=
type ModInfo struct {
	DisplayName        string `json:"display_name"`
	Rank               int    `json:"rank"`
	InternalName       string `json:"internal_name"`
	Author             string `json:"author"`
	Homepage           string `json:"homepage"`
	Description        string `json:"description"`
	Icon               string `json:"icon"`
	Version            string `json:"version"`
	TModLoaderVersion  string `json:"tmodloader_version"`
	LastUpdated        string `json:"last_updated"`
	ModDependencies    string // added later
	ModSide            string `json:"modside"`
	DownloadLink       string `json:"download_link"`
	DownloadsTotal     int    `json:"downloads_total"`
	DownloadsToday     int    `json:"downloads_today"`
	DownloadsYesterday int    `json:"downloads_yesterday"`
}

type AuthorModInfo struct {
	Rank               int    `json:"rank"`
	DisplayName        string `json:"display_name"`
	EscapedDisplayName string // added later
	DownloadsTotal     int    `json:"downloads_total"`
	DownloadsYesterday int    `json:"downloads_yesterday"`
}

type AuthorMaintainedModInfo struct {
	InternalName       string `json:"internal_name"`
	EscapedModName     string // added later
	DownloadsTotal     int    `json:"downloads_total"`
	DownloadsYesterday int    `json:"downloads_yesterday"`
}

type Author struct {
	SteamName          string                    `json:"steam_name"`
	DownloadsTotal     int                       `json:"downloads_total"`
	DownloadsYesterday int                       `json:"downloads_yesterday"`
	Mods               []AuthorModInfo           `json:"mods"`
	MaintainedMods     []AuthorMaintainedModInfo `json:"maintained_mods"`
}

var ModNameMap map[string]string = make(map[string]string)             //maps Display names to Internal names
var ModInfoMap map[string]*ModListItem = make(map[string]*ModListItem) //maps Internal names to ModListItem (for Rank and DownloadsToday)
var ModList []ModListItem = make([]ModListItem, 0)

/*!!!IMPORTANT always lock the mutex below before working with the data above, and close it afterwards!!!*/
var dataMutex *sync.Mutex = &sync.Mutex{} //while the maps and slices are being updated or used this mutex will be locked

func updateModMaps() error {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	resp, err := http.Get("https://tmlapis.tomat.dev/1.3/list")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var TempModList []ModListItem = make([]ModListItem, len(ModList)) //if the fetching fails, it is not a fatal error as we still have the old modlist
	err = json.NewDecoder(resp.Body).Decode(&TempModList)             //decode the modlist
	if err != nil {
		return err
	}
	ModList = TempModList

	ModNameMap = make(map[string]string)
	for i := range ModList {
		ModNameMap[url.QueryEscape(ModList[i].DisplayName)] = ModList[i].InternalName //map all Display names to Internal names
		ModInfoMap[ModList[i].InternalName] = &ModList[i]                             //map all Internal names to ModInfo data
	}
	sort.Slice(ModList, func(i, j int) bool {
		return ModList[i].Rank < ModList[j].Rank
	})
	return nil
}
