package tmlapi14

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/NotLe0n/ModStats/server/config"
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

type ModVersionData struct {
	ModVersion        string `json:"mod_version"`
	TmodloaderVersion string `json:"tmodloader_version"`
}

// ModInfo holds mod info that is fetched from tmlapis.tomat.dev/1.4/mod
type ModInfo struct {
	DisplayName        string           `json:"display_name"`
	InternalName       string           `json:"internal_name"`
	ModID              string           `json:"mod_id"`
	Author             string           `json:"author"`
	AuthorID           string           `json:"author_id"`
	ModSide            string           `json:"modside"`
	Homepage           string           `json:"homepage"`
	Versions           []ModVersionData `json:"versions"`
	ModReferences      string           `json:"mod_references"`
	NumVersions        uint32           `json:"num_versions"`
	Tags               []ModTag         `json:"tags"`
	TimeCreated        uint64           `json:"time_created"`
	TimeUpdated        uint64           `json:"time_updated"`
	IconUrl            string           `json:"workshop_icon_url"`
	DownloadsTotal     uint32           `json:"downloads_total"`
	Favorited          uint32           `json:"favorited"`
	Followers          uint32           `json:"followers"`
	Views              uint64           `json:"views"`
	VoteData           *ModVoteData     `json:"vote_data"`
	Playtime           string           `json:"playtime"`
	DisplayNameHTML    template.HTML    // added later
	Description        string           `json:"description"`
	EscapedDescription template.HTML    // added later
	Children           []ModChild       `json:"children"`
}

type Author struct {
	SteamID        uint64    `json:"steam_id"`
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
	return modList
}

func updateModMaps() error {
	var TempmodList = make([]ModInfo, len(modList)) //if the fetching fails, it is not a fatal error as we still have the old modList

	used_data := false
	use_data := func(err error) error {
		if len(modList) == 0 {
			logf("error loading mod list: %s\nno state available, loading data/modlist_1_4.json", err.Error())
			// load data/modlist_1_3.json without using helper
			file, err := os.Open("data/modlist_1_4.json")
			if err != nil {
				return err
			}
			defer file.Close()

			used_data = true
			// decode the data into a temp list
			return json.NewDecoder(file).Decode(&TempmodList)
		}
		return err
	}

	// get the data
	resp, err := helper.GetWithTimeout(config.C.GetString("API-URL") + "/1.4/list")
	if err != nil {
		if err := use_data(err); err != nil {
			return err
		}
	} else {
		defer resp.Body.Close()
	}

	if !used_data {
		// decode the data into a temp list
		if err := json.NewDecoder(resp.Body).Decode(&TempmodList); err != nil {
			if err := use_data(err); err != nil {
				return err
			}
		}
	}

	bbcodeCompiler := helper.NewBBCodeToTextCompiler()
	for i := range TempmodList {
		TempmodList[i].DisplayNameHTML = template.HTML(helper.ParseChatTags(TempmodList[i].DisplayName))

		parsedBB := bbcodeCompiler.Compile(TempmodList[i].Description)
		unesc := helper.Unescape(parsedBB)
		TempmodList[i].EscapedDescription = template.HTML(strings.TrimLeft(unesc, "<br> "))
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
	logf("checking '%s/1.4/'...", config.C.GetString("API-URL"))
	if _, err := helper.GetWithTimeout(config.C.GetString("API-URL") + "/1.4/"); err != nil {
		logf("'%s/1.4/)' can't be reached, switching to secondary mirror.", config.C.GetString("API-URL"))
	}

	// adds logging
	intitUpdate := func() {
		logf("updating ModNameMap14")
		if err := updateModMaps(); err != nil {
			logf("Unable to update ModNameMap14, using the last valid state: %s", err.Error())
		} else {
			logf("done updating ModNameMap14")
		}
	}

	// update once at the beginning
	intitUpdate()

	//this goroutine updates the mod list every 15 minutes so that the loading time is not too long on every reload
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			intitUpdate()
		}
	}()
}

func logf(format string, args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, "[LOG] "+format+"\n", args...)
}
