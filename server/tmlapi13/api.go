package tmlapi13

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/NotLe0n/ModStats/server/helper"
)

var apiUrl = "https://tmlapis.repl.co/1.3/"

// requires a internal modName
func GetModInfo(modName string) (ModInfo, error) {
	resp, err := http.Get(apiUrl + "mod/" + url.QueryEscape(modName))
	if err != nil {
		return ModInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ModInfo{}, fmt.Errorf("tmlapis returned %d", resp.StatusCode)
	}

	var modInfo ModInfo
	if err := helper.UnmarshalResp(resp.Body, &modInfo); err != nil { //encode the data (without rank and DownloadsToday)
		return ModInfo{}, err
	}

	dataMutex.RLock()
	defer dataMutex.RUnlock()

	modInfo.Rank = modInfoMap[modInfo.InternalName].Rank                     //add the rank from the map
	modInfo.DownloadsToday = modInfoMap[modInfo.InternalName].DownloadsToday //add DownloadsToday from the map
	modInfo.TModLoaderVersion = modInfo.TModLoaderVersion[11:]

	return modInfo, nil
}

type ModVersion struct {
	Version           string `json:"version"`
	DownloadsTotal    int    `json:"downloads_total"`
	TModLoaderVersion string `json:"tmodloader_version"`
	PublishDate       string `json:"publish_date"`
}

// requires a internal modName
func GetModVersionHistory(modName string) ([]ModVersion, error) {
	resp, err := http.Get(apiUrl + "history/" + url.QueryEscape(modName)) //fetch most of the data
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tmlapis returned %d", resp.StatusCode)
	}

	var modVersion []ModVersion
	if err := helper.UnmarshalResp(resp.Body, &modVersion); err != nil {
		return nil, err
	}

	return modVersion, err
}

func GetAuthorInfo(steamid64 string) (Author, error) {
	resp, err := http.Get(apiUrl + "author/" + steamid64) //fetch data
	if err != nil {
		return Author{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Author{}, fmt.Errorf("tmlapis returned %d", resp.StatusCode)
	}

	var authorInfo Author
	if err := helper.UnmarshalResp(resp.Body, &authorInfo); err != nil {
		return Author{}, err
	}

	dataMutex.RLock()
	defer dataMutex.RUnlock()

	for i := range authorInfo.Mods {
		authorInfo.Mods[i].InternalName = url.QueryEscape(modNameMap[url.QueryEscape(authorInfo.Mods[i].DisplayName)])
	}
	for i := range authorInfo.MaintainedMods {
		authorInfo.MaintainedMods[i].EscapedModName = url.QueryEscape(authorInfo.MaintainedMods[i].InternalName)
	}

	return authorInfo, nil
}

func init() {
	rand.Seed(time.Now().Unix())
}

// returns a random internal mod name
func GetRandomMod() string {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	n := rand.Intn(len(modList))
	return modList[n].InternalName
}
