package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func logf(format string, args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, "[LOG] "+format+"\n", args...)
}

func main() {

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
