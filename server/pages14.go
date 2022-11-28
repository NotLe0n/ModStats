package main

import (
	"fmt"
	"github.com/NotLe0n/ModStats/server/helper"
	"github.com/NotLe0n/ModStats/server/tmlapi13"
	"github.com/NotLe0n/ModStats/server/tmlapi14"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func indexPage14(c *gin.Context) {
	ModList := tmlapi14.GetModList()

	combinedDownloads := func(modList []tmlapi14.ModInfo) (combined int) {
		for i := range modList {
			combined += int(modList[i].DownloadsTotal)
		}
		return combined
	}

	c.HTML(http.StatusOK, "base/index.gohtml", gin.H{
		"modlist":   ModList,
		"modcount":  len(ModList),
		"combined":  combinedDownloads(ModList),
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
		c.AbortWithStatus(http.StatusBadRequest)
	}

	authorInfo, err := tmlapi14.GetAuthorInfo(authorID)
	if err != nil {
		logf("Error getting authorInfo: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
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
		c.AbortWithStatus(http.StatusBadRequest)
	}

	modData, err := tmlapi14.GetModInfo(modName)
	if err != nil {
		logf("Error getting modInfo for %s: %s", modName, err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
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
