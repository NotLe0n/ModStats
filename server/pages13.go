package main

import (
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/NotLe0n/ModStats/server/helper"
	"github.com/NotLe0n/ModStats/server/tmlapi13"
	"github.com/gin-gonic/gin"
)

func indexPage13(c *gin.Context) {
	ModList := tmlapi13.GetModList()

	combinedDownloads := func(modList []tmlapi13.ModListItem) (combined int) {
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

	hotMods := func() []tmlapi13.ModListItem {
		hotMods := make([]tmlapi13.ModListItem, 10)
		sort.Slice(ModList, func(i, j int) bool {
			return ModList[i].DownloadsYesterday > ModList[j].DownloadsYesterday
		})
		copy(hotMods, ModList)
		sort.Slice(ModList, func(i, j int) bool {
			return ModList[i].Rank < ModList[j].Rank
		})
		return hotMods
	}
	percent := 0.0
	if len(ModList) >= 10 {
		percent = float64(combinedDownloads(ModList[:10])) / float64(combinedDownloads(ModList)) * 100
	}
	c.HTML(http.StatusOK, "base/index.gohtml", gin.H{
		"modlist":   ModList,
		"modcount":  len(ModList),
		"combined":  strconv.FormatFloat(float64(combinedDownloads(ModList))/1_000_000.0, 'f', 3, 64),
		"deadmods":  deadMods(),
		"percent":   strconv.FormatFloat(percent, 'f', 2, 64),
		"median":    ModList[len(ModList)/2].DownloadsTotal,
		"contribs":  80,
		"top10mods": ModList[:10],
		"hotmods":   hotMods(),
		"isLegacy":  true,
	})
}

func modListPage13(c *gin.Context) {
	c.HTML(http.StatusOK, "base/list.gohtml", gin.H{
		"modlist":  tmlapi13.GetModList(),
		"isLegacy": true,
	})
}

func authorPage13(c *gin.Context) {
	authorID := c.Param("authorID")
	if authorID == "" {
		c.HTML(http.StatusBadRequest, "base/error.gohtml", gin.H{
			"modlist":  tmlapi13.GetModList(),
			"error":    "Please search for a author's name or SteamID.",
			"isLegacy": true,
		})
		return
	}

	authorInfo, err := tmlapi13.GetAuthorInfo(authorID)
	if err != nil {
		logf("Error getting authorInfo: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "base/error.gohtml", gin.H{
			"modlist":  tmlapi13.GetModList(),
			"error":    "An author with the name or SteamID '" + authorID + "' was not found or a internal error has occurred.",
			"isLegacy": true,
		})
		return
	}

	for i, mod := range authorInfo.Mods {
		authorInfo.Mods[i].DisplayName = helper.ParseChatTags(helper.Unescape(mod.DisplayName))
	}

	c.HTML(http.StatusOK, "base/author.gohtml", gin.H{
		"modlist":    tmlapi13.GetModList(),
		"author":     authorID,
		"authorInfo": authorInfo,
		"isLegacy":   true,
	})
}

func modStatsPage13(c *gin.Context) {
	modName := c.Param("modID")
	if modName == "" {
		c.HTML(http.StatusBadRequest, "base/error.gohtml", gin.H{
			"modlist":  tmlapi13.GetModList(),
			"error":    "Please search for a mod's name.",
			"isLegacy": true,
		})
		return
	}

	if internal_name, ok := tmlapi13.GetInternalName(url.QueryEscape(modName)); ok {
		modName = internal_name
	}

	modData, err := tmlapi13.GetModInfo(modName)
	if err != nil {
		logf("Error getting modInfo for %s: %s", modName, err.Error())
		c.HTML(http.StatusInternalServerError, "base/error.gohtml", gin.H{
			"modlist":  tmlapi13.GetModList(),
			"error":    "A mod with the name '" + modName + "' was not found or a internal error has occurred.",
			"isLegacy": true,
		})
		return
	}

	if resp, err := http.Get(modData.Icon); resp.StatusCode == http.StatusNotFound || err != nil {
		modData.Icon = ""
	}

	modDependencies := strings.Split(modData.ModDependencies, ", ")
	if modDependencies[0] == "" {
		modDependencies = make([]string, 0)
	}

	modVersions, err := tmlapi13.GetModVersionHistory(modName)
	if err != nil {
		logf("Error getting modVersionHistory: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "base/mod.gohtml", gin.H{
		"modlist":            tmlapi13.GetModList(),
		"modData":            modData,
		"modDependencies":    modDependencies,
		"versionHistory":     modVersions,
		"escapedDisplayName": template.HTML(helper.ParseChatTags(helper.Unescape(modData.DisplayName))),
		"escapedDescription": template.HTML(helper.ParseChatTags(helper.Unescape(modData.Description))),
		"isLegacy":           true,
	})
}
