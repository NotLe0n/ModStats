package main

import (
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

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

	c.HTML(http.StatusOK, "index.gohtml", gin.H{
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

	c.HTML(http.StatusOK, "list.gohtml", gin.H{
		"modlist": ModList,
	})
}

func authorStatsPage(c *gin.Context) {
	authorID := c.Param("authorID")
	if authorID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	authorInfo, err := getAuthorInfo(authorID)
	if err != nil {
		logf("Error getting authorInfo: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "author.gohtml", gin.H{
		"modlist":    ModList,
		"author":     authorID,
		"authorInfo": authorInfo,
	})
}

func modStatsPage(c *gin.Context) {
	modName := c.Param("modID")
	if modName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	isLegacy := c.Request.URL.Query().Has("legacy")

	dataMutex.Lock()
	if name, ok := ModNameMap[url.QueryEscape(modName)]; ok {
		modName = name
	}
	dataMutex.Unlock()

	modData, err := getModInfo(modName)
	if err != nil {
		logf("Error getting modInfo for %s: %s", modName, err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if resp, err := http.Get(modData.Icon); resp.StatusCode == 404 || err != nil {
		modData.Icon = ""
	}

	modDependencies := strings.Split(modData.ModDependencies, ", ")
	if modDependencies[0] == "" {
		modDependencies = make([]string, 0)
	}

	modVersions, err := getModVersionHistory(modName)
	if err != nil {
		logf("Error getting modVersionHistory: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dataMutex.Lock()
	defer dataMutex.Unlock()

	c.HTML(http.StatusOK, "mod.gohtml", gin.H{
		"isLegacy":           isLegacy,
		"modlist":            ModList,
		"modData":            modData,
		"modDependencies":    modDependencies,
		"versionHistory":     modVersions,
		"escapedDisplayName": template.HTML(ParseChatTags(modData.DisplayName)),
		"escapedDescription": template.HTML(ParseChatTags(strings.Trim(modData.Description, "\""))),
	})
}
