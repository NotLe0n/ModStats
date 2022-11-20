package tmlapi13

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/NotLe0n/ModStats/server/helper"

	"github.com/gin-gonic/gin"
)

// ModListItem holds a single item when fetching the whole modList
type ModListItem struct {
	DisplayName        string        `json:"display_name"`
	InternalName       string        `json:"internal_name"`
	DisplayNameHTML    template.HTML // added later
	Rank               int           `json:"rank"`
	DownloadsTotal     int           `json:"downloads_total"`
	DownloadsToday     int           `json:"downloads_today"`
	DownloadsYesterday int           `json:"downloads_yesterday"`
	TModLoaderVersion  string        `json:"tmodloader_version"`
}

// ModInfo holds mod info that is fetched from tmlapis.tomat.dev/1.3/mod
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
	ModDependencies    string `json:"modreferences"`
	ModSide            string `json:"modside"`
	DownloadLink       string `json:"download_link"`
	DownloadsTotal     int    `json:"downloads_total"`
	DownloadsToday     int    `json:"downloads_today"`
	DownloadsYesterday int    `json:"downloads_yesterday"`
}

type AuthorModInfo struct {
	Rank               int    `json:"rank"`
	DisplayName        string `json:"display_name"`
	InternalName       string // added later
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

var modNameMap = make(map[string]string)       //maps Display names to Internal names
var modInfoMap = make(map[string]*ModListItem) //maps Internal names to ModListItem (for Rank and DownloadsToday)
var modList = make([]ModListItem, 0)

/*!!!IMPORTANT always lock the mutex below before working with the data above, and close it afterwards!!!*/
var dataMutex = &sync.RWMutex{} //while the maps and slices are being updated or used this mutex will be locked

func GetInternalName(display_name string) (string, bool) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	internal_name, ok := modNameMap[display_name]
	return internal_name, ok
}

func GetModList() []ModListItem {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	copyList := make([]ModListItem, len(modList))
	copy(copyList, modList)
	return copyList
}

func updateModMaps() error {
	// get the data
	resp, err := http.Get("https://tmlapis.tomat.dev/1.3/list")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// decode the data into a temp list
	var TempmodList = make([]ModListItem, len(modList))   //if the fetching fails, it is not a fatal error as we still have the old modList
	err = json.NewDecoder(resp.Body).Decode(&TempmodList) //decode the modList
	if err != nil {
		return err
	}

	// lock the mutex for writing
	dataMutex.Lock()
	defer dataMutex.Unlock()

	modList = TempmodList

	modNameMap = make(map[string]string)
	for i := range modList {
		modList[i].DisplayNameHTML = template.HTML(helper.ParseChatTags(modList[i].DisplayName))
		modNameMap[url.QueryEscape(modList[i].DisplayName)] = modList[i].InternalName //map all Display names to Internal names
		modInfoMap[modList[i].InternalName] = &modList[i]                             //map all Internal names to ModInfo data
	}

	sort.Slice(modList, func(i, j int) bool {
		return modList[i].Rank < modList[j].Rank
	})

	return nil
}

// start the ticker to update the state
func init() {
	//this goroutine updates the mod list every 15 minutes so that the loading time is not too long on every reload

	go func() {
		for ticker := time.NewTicker(15 * time.Minute); true; <-ticker.C {
			logf("updating ModNameMap")
			if err := updateModMaps(); err != nil {
				logf("Unable to update ModNameMap, using the last valid state: %s", err.Error())
			}
		}
	}()
}

func logf(format string, args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, "[LOG] "+format+"\n", args...)
}