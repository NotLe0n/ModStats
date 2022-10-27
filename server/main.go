package main

import (
	"fmt"
	"log"
	"math/rand"
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
	r.GET("/modList", modListPage)
	r.GET("/stats", statsPage)
	api := r.Group("/api")
	{
		api.GET("/getRandomMod", getRandomMod)
	}

	log.Fatal(r.Run())
}
