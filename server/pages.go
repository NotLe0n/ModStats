package main

import (
	"html/template"
	"net/http"
	"net/url"
	"regexp"
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

	c.HTML(http.StatusOK, "modList.gohtml", gin.H{
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

	iconDisplay := "block"
	if modData.Icon == "" {
		iconDisplay = "none"
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

	parseChatTags := func(str string) string {
		rpa := strings.ReplaceAll
		str = rpa(
			rpa(
				rpa(
					rpa(
						rpa(
							rpa(
								rpa(
									rpa(
										str,
										"<",
										"&lt;",
									),
									">",
									"&gt;",
								),
								"\r\n",
								"<br>",
							),
							"\n",
							"<br>",
						),
						"\t",
						"    ",
					),
					"\\'",
					"'",
				),
				"\\\\",
				"\\",
			),
			"\\\"",
			"&quot",
		)
		itemTagRegex := regexp.MustCompile("\\[i(.*?):(\\w+)\\]")
		str = itemTagRegex.ReplaceAllString(str, "<img src=\"https://tmlapis.repl.co/img/Item_$2.png\" id=\"item-icon\">")
		colorTagRegex := regexp.MustCompile("\\[c\\/(\\w+):([\\s\\S]+?)\\]")
		str = colorTagRegex.ReplaceAllString(str, "<span style=\"color: #$1;\">$2</span>")
		return str
	}

	dataMutex.Lock()
	defer dataMutex.Unlock()
	c.HTML(http.StatusOK, "mod.gohtml", gin.H{
		"modlist":            ModList,
		"modData":            modData,
		"iconDisplay":        iconDisplay,
		"hasHomepage":        modData.Homepage != "",
		"modDependencies":    modDependencies,
		"versionHistory":     modVersions,
		"escapedDisplayName": template.HTML(parseChatTags(modData.DisplayName)),
		"escapedDescription": template.HTML(parseChatTags(strings.Trim(modData.Description, "\""))),
	})
}
