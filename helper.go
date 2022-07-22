package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// requires a internal modName
func getModInfo(modName string) (ModInfo, error) {
	resp, err := http.Get("https://tmlapis.repl.co/modInfo?modname=" + url.QueryEscape(modName))
	if err != nil {
		return ModInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ModInfo{}, fmt.Errorf("tmlapis returned %d", resp.StatusCode)
	}

	var modInfo ModInfo
	err = json.NewDecoder(resp.Body).Decode(&modInfo) //encode the data (without rank and DownloadsToday)
	if err != nil {
		return ModInfo{}, err
	}
	dataMutex.Lock()
	defer dataMutex.Unlock()
	modInfo.Rank = ModInfoMap[modInfo.InternalName].Rank                     //add the rank from the map
	modInfo.DownloadsToday = ModInfoMap[modInfo.InternalName].DownloadsToday //add DownloadsToday from the map
	return modInfo, nil
}

type ModVersion struct {
	Version           string
	Downloads         int
	TModLoaderVersion string
	PublishDate       string
}

// requires a internal modName
func getModVersionHistory(modName string) ([]ModVersion, error) {
	resp, err := http.Get("https://tmlapis.repl.co/modVersionHistory?modname=" + url.QueryEscape(modName)) //fetch most of the data
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmlapis returned %d", resp.StatusCode)
	}

	var modInfo []ModVersion
	err = json.NewDecoder(resp.Body).Decode(&modInfo) //encode the data (without rank and DownloadsToday)
	if err != nil {
		return nil, err
	}
	return modInfo, err
}

func getAuthorInfo(steamid64 string) (Author, error) {
	resp, err := http.Get("https://tmlapis.repl.co/author_api/" + steamid64) //fetch data
	if err != nil {
		return Author{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Author{}, fmt.Errorf("tmlapis returned %d", resp.StatusCode)
	}

	var authorInfo Author
	err = json.NewDecoder(resp.Body).Decode(&authorInfo) //encode the data (without rank and DownloadsToday)
	if err != nil {
		return Author{}, err
	}
	for i := range authorInfo.Mods {
		authorInfo.Mods[i].EscapedDisplayName = url.QueryEscape(authorInfo.Mods[i].DisplayName)
	}
	for i := range authorInfo.MaintainedMods {
		authorInfo.MaintainedMods[i].EscapedModName = url.QueryEscape(authorInfo.MaintainedMods[i].ModName)
	}
	return authorInfo, nil
}
