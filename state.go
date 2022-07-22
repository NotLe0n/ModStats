package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"sync"
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

type AuthorModInfo struct {
	RankTotal          int
	DisplayName        string
	EscapedDisplayName string
	DownloadsTotal     int
	DownloadsYesterday int
}

type AuthorMaintainedModInfo struct {
	ModName            string
	EscapedModName     string
	DownloadsTotal     int
	DownloadsYesterday int
}

type Author struct {
	SteamName          string
	DownloadsTotal     int
	DownloadsYesterday int
	Mods               []AuthorModInfo
	MaintainedMods     []AuthorMaintainedModInfo
}

var ModNameMap map[string]string = make(map[string]string)             //maps Display names to Internal names
var ModInfoMap map[string]*ModListItem = make(map[string]*ModListItem) //maps Internal names to ModListItem (for Rank and DownloadsToday)
var ModList []ModListItem = make([]ModListItem, 0)

/*!!!IMPORTANT always lock the mutex below before working with the data above, and close it afterwards!!!*/
var dataMutex *sync.Mutex = &sync.Mutex{} //while the maps and slices are being updated or used this mutex will be locked

func updateModMaps() error {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	resp, err := http.Get("https://tmlapis.repl.co/modList")
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
		ModNameMap[url.QueryEscape(ModList[i].DisplayName)] = ModList[i].ModName //map all Display names to Internal names
		ModInfoMap[ModList[i].ModName] = &ModList[i]                             //map all Internal names to ModInfo data
	}
	sort.Slice(ModList, func(i, j int) bool {
		return ModList[i].Rank < ModList[j].Rank
	})
	return nil
}
