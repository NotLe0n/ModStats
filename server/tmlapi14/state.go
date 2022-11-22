package tmlapi14

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/NotLe0n/ModStats/server/helper"

	"github.com/gin-gonic/gin"
)

type ModChild struct {
	File_type       uint32 `json:"file_type"`
	Publishedfileid string `json:"publishedfileid"`
	Sortorder       uint32 `json:"sortorder"`
}

type ModTag struct {
	Tag         string `json:"tag"`
	DisplayName string `json:"display_name"`
}

type ModVoteData struct {
	Score     float64 `json:"score"`
	VotesUp   uint32  `json:"votes_up"`
	VotesDown uint32  `json:"votes_down"`
}

// ModInfo holds mod info that is fetched from tmlapis.tomat.dev/1.4/mod
type ModInfo struct {
	DisplayName       string        `json:"display_name"`
	InternalName      string        `json:"internal_name"`
	ModID             string        `json:"mod_id"`
	Author            string        `json:"author"`
	AuthorID          string        `json:"author_id"`
	ModSide           string        `json:"modside"`
	Homepage          string        `json:"homepage"`
	TModLoaderVersion string        `json:"tmodloader_version"`
	Version           string        `json:"version"`
	ModReferences     string        `json:"mod_references"`
	NumVersions       uint32        `json:"num_versions"`
	Tags              []ModTag      `json:"tags"`
	TimeCreated       uint64        `json:"time_created"`
	TimeUpdated       uint64        `json:"time_updated"`
	IconUrl           string        `json:"workshop_icon_url"`
	DownloadsTotal    uint32        `json:"downloads_total"`
	Favorited         uint32        `json:"favorited"`
	Followers         uint32        `json:"followers"`
	Views             uint64        `json:"views"`
	VoteData          *ModVoteData  `json:"vote_data"`
	Playtime          string        `json:"playtime"`
	DisplayNameHTML   template.HTML // added later
	Description       string        `json:"description"`
	Children          []ModChild    `json:"children"`
}

type Author struct {
	SteamName      string    `json:"steam_name"`
	Mods           []ModInfo `json:"mods"`
	Total          uint32    `json:"total"`
	TotalDownloads uint64    `json:"total_downloads"`
	TotalFavorites uint64    `json:"total_favorites"`
	TotalViews     uint64    `json:"total_views"`
}

var modList = make([]ModInfo, 0)

/*!!!IMPORTANT always lock the mutex below before working with the data above, and close it afterwards!!!*/
var dataMutex = &sync.RWMutex{} //while the maps and slices are being updated or used this mutex will be locked

func GetModList() []ModInfo {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	copyList := make([]ModInfo, len(modList))
	copy(copyList, modList)
	return copyList
}

func updateModMaps() error {
	// get the data
	resp, err := http.Get("https://tmlapis.tomat.dev/1.4/list")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// decode the data into a temp list
	var TempmodList = make([]ModInfo, len(modList))       //if the fetching fails, it is not a fatal error as we still have the old modList
	err = json.NewDecoder(resp.Body).Decode(&TempmodList) //decode the modList
	if err != nil {
		return err
	}

	for i := range TempmodList {
		TempmodList[i].DisplayNameHTML = template.HTML(helper.ParseChatTags(TempmodList[i].DisplayName))
	}

	sort.Slice(TempmodList, func(i, j int) bool {
		return TempmodList[i].DownloadsTotal > TempmodList[j].DownloadsTotal
	})

	// lock the mutex for writing
	dataMutex.Lock()
	defer dataMutex.Unlock()

	modList = TempmodList

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
