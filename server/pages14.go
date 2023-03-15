package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/NotLe0n/ModStats/server/helper"
	"github.com/NotLe0n/ModStats/server/tmlapi13"
	"github.com/NotLe0n/ModStats/server/tmlapi14"
	"github.com/gin-gonic/gin"
)

func indexPage14(c *gin.Context) {
	ModList := tmlapi14.GetModList()

	if len(ModList) < 10 {
		c.HTML(http.StatusInternalServerError, "base/error.gohtml", gin.H{
			"modlist": tmlapi13.GetModList(),
			"error":   "unable to fetch valid Mod List",
		})
		return
	}

	combinedDownloads := func(modList []tmlapi14.ModInfo) (combined int) {
		for i := range modList {
			combined += int(modList[i].DownloadsTotal)
		}
		return combined
	}

	c.HTML(http.StatusOK, "base/index.gohtml", gin.H{
		"modlist":   ModList,
		"modcount":  len(ModList),
		"combined":  strconv.FormatFloat(float64(combinedDownloads(ModList))/1_000_000.0, 'f', 3, 64),
		"percent":   strconv.FormatFloat(float64(combinedDownloads(ModList[:10]))/float64(combinedDownloads(ModList))*100, 'f', 2, 64),
		"median":    ModList[len(ModList)/2].DownloadsTotal,
		"contribs":  150,
		"top10mods": ModList[:10],
		"isLegacy":  false,
	})
}

func modListPage14(c *gin.Context) {
	c.HTML(http.StatusOK, "base/list.gohtml", gin.H{
		"modlist":  tmlapi14.GetModList(),
		"isLegacy": false,
	})
}

func authorPage14(c *gin.Context) {
	authorID := c.Param("authorID")
	if authorID == "" {
		c.HTML(http.StatusBadRequest, "base/error.gohtml", gin.H{
			"modlist":  tmlapi13.GetModList(),
			"error":    "Please search for a author's name or SteamID.",
			"isLegacy": true,
		})
		return
	}

	authorInfo, err := tmlapi14.GetAuthorInfo(authorID)
	if err != nil {
		logf("Error getting authorInfo: %s", err.Error())
		c.HTML(http.StatusInternalServerError, "base/error.gohtml", gin.H{
			"modlist":  tmlapi13.GetModList(),
			"error":    "An author with the name or SteamID " + authorID + " was not found or a internal error has occurred.",
			"isLegacy": true,
		})
		return
	}

	c.HTML(http.StatusOK, "base/author.gohtml", gin.H{
		"modlist":    tmlapi14.GetModList(),
		"authorID":   authorInfo.SteamID,
		"authorInfo": authorInfo,
		"isLegacy":   false,
	})
}

func modStatsPage14(c *gin.Context) {
	modName := c.Param("modID")
	if modName == "" {
		c.HTML(http.StatusBadRequest, "base/error.gohtml", gin.H{
			"modlist":  tmlapi13.GetModList(),
			"error":    "Please search for a mod's name or file id.",
			"isLegacy": true,
		})
		return
	}

	modData, err := tmlapi14.GetModInfo(modName)
	if err != nil {
		logf("Error getting modInfo for %s: %s", modName, err.Error())
		c.HTML(http.StatusInternalServerError, "base/error.gohtml", gin.H{
			"modlist":  tmlapi13.GetModList(),
			"error":    "A mod with the name or file id '" + modName + "' was not found or a internal error has occurred.",
			"isLegacy": true,
		})
		return
	}

	if resp, err := http.Get(modData.IconUrl); resp.StatusCode == http.StatusNotFound || err != nil {
		modData.IconUrl = ""
	}

	var a, b, cc, d int
	_, _ = fmt.Sscanf(strings.ReplaceAll(modData.TModLoaderVersion, "tModLoader v", ""), "%d.%d.%d.%d", &a, &b, &cc, &d)
	modData.TModLoaderVersion = fmt.Sprintf("v%d.%02d.%d.%d", a, b, cc, d)
	modData.DisplayNameHTML = template.HTML(helper.ParseChatTags(modData.DisplayName))

	for i, tag := range modData.Tags {
		modData.Tags[i].DisplayName = strings.Title(tag.DisplayName)
	}

	numFilledStars := int(5 * modData.VoteData.Score)
	numEmptyStars := 5 - numFilledStars
	stars := ""

	for i := 0; i < numFilledStars; i++ {
		stars += "★"
	}
	for i := 0; i < numEmptyStars; i++ {
		stars += "☆"
	}

	bbcodeCompiler := helper.NewBBCodeCompiler()

	c.HTML(http.StatusOK, "base/mod.gohtml", gin.H{
		"modlist":            tmlapi13.GetModList(),
		"modData":            modData,
		"escapedDisplayName": template.HTML(helper.ParseChatTags(modData.DisplayName)),
		"escapedDescription": template.HTML(helper.ParseChatTags(bbcodeCompiler.Compile(modData.Description))),
		"stars":              stars,
		"isLegacy":           false,
	})
}
