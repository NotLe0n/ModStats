package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func logf(format string, args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, "[LOG] "+format+"\n", args...)
}

func startTicker() {
	//this goroutine updates the mod list every 10 minutes so that the loading time is not too long on every reload
	go func() {
		for ; true; <-time.Tick(15 * time.Minute) {
			logf("updating ModNameMap")
			if err := updateModMaps(); err != nil {
				logf("Unable to update ModNameMap, using the last valid state: %s", err.Error())
			}
		}
	}()
}

func main() {
	rand.Seed(time.Now().Unix())
	startTicker()

	r := gin.Default()

	r.LoadHTMLGlob("./html/*")

	r.Static("/static", "./static")
	r.Static("/favicon.ico", "./favicon.ico")

	r.GET("/", indexPage)
	r.GET("/list", modListPage)
	r.GET("/mod/:modID", modStatsPage)
	r.GET("/author/:authorID", authorStatsPage)
	api := r.Group("/api")
	{
		api.GET("/getRandomMod", getRandomMod)
	}

	log.Fatal(r.Run())
}

func ParseChatTags(str string) string {
	replaceMap := map[string]string{
		"<":    "&lt;",
		">":    "&gt;",
		"\r\n": "<br>",
		"\n":   "<br>",
		"\t":   "    ",
		"\\'":  "'",
		"\\\\": "\\",
		"\\\"": "&quot",
	}

	for oldStr, newStr := range replaceMap {
		str = strings.ReplaceAll(str, oldStr, newStr)
	}

	itemTagRegex := regexp.MustCompile("\\[i(.*?):(\\w+)\\]")
	str = itemTagRegex.ReplaceAllString(str, "<img src=\"https://tmlapis.repl.co/img/Item_$2.png\" id=\"item-icon\">")
	colorTagRegex := regexp.MustCompile("\\[c\\/(\\w+):([\\s\\S]+?)\\]")
	str = colorTagRegex.ReplaceAllString(str, "<span style=\"color: #$1;\">$2</span>")
	return str
}
