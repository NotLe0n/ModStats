package tmlapi14

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/NotLe0n/ModStats/server/config"
	"github.com/NotLe0n/ModStats/server/helper"
)

// requires a internal modName
func GetModInfo(modName string) (ModInfo, error) {
	resp, err := helper.GetWithTimeout(config.C.GetString("API-URL") + "/1.4/mod/" + url.QueryEscape(modName))
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

	return modInfo, nil
}

func GetAuthorInfo(steamid64 string) (Author, error) {
	resp, err := helper.GetWithTimeout(config.C.GetString("API-URL") + "/1.4/author/" + steamid64) //fetch data
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
