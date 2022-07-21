package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func indexPage(c *gin.Context) {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	combinedDownloads := func(modList []ModListItem) (combined int) {
		for i := range modList {
			combined += modList[i].DownloadsTotal
		}
		return combined
	}

	deadMods := func() (deadmods int) {
		for i := range ModList {
			if ModList[i].DownloadsYesterday < 5 {
				deadmods++
			}
		}
		return deadmods
	}

	hotMods := func() []ModListItem {
		hotMods := make([]ModListItem, 10)
		sort.Slice(ModList, func(i, j int) bool {
			return ModList[i].DownloadsYesterday > ModList[j].DownloadsYesterday
		})
		copy(hotMods, ModList)
		sort.Slice(ModList, func(i, j int) bool {
			return ModList[i].Rank < ModList[j].Rank
		})
		return hotMods
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"modlist":   ModList,
		"modcount":  len(ModList),
		"combined":  combinedDownloads(ModList),
		"deadmods":  deadMods(),
		"percent":   strconv.FormatFloat(float64(combinedDownloads(ModList[:10]))/float64(combinedDownloads(ModList))*100, 'f', 2, 64),
		"median":    ModList[len(ModList)/2].DownloadsTotal,
		"contribs":  80,
		"top10mods": ModList[:10],
		"hotmods":   hotMods(),
	})
}

func modListPage(c *gin.Context) {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	c.HTML(http.StatusOK, "modList.html", gin.H{
		"modlist": ModList,
	})
}

func statsPage(c *gin.Context) {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	steamid64 := c.Query("author")
	if steamid64 != "" {
		c.String(http.StatusInternalServerError, "not implemented")
		return
	}

	modName := c.Query("mod")
	if modName != "" {
		c.HTML(http.StatusOK, "stats.html", gin.H{
			"Mod": modName,
		})
		return
	}

	c.AbortWithStatus(http.StatusBadRequest)
}
