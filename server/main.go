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

	r.LoadHTMLGlob("./html/**/*.gohtml")

	r.Static("/static", "./static")
	r.Static("/favicon.ico", "./static/assets/favicon.ico")

	legacy := r.Group("/legacy")
	{
		legacy.GET("/", indexPage13)
		legacy.GET("/list", modListPage13)
		legacy.GET("/mod/", modStatsPage13)
		legacy.GET("/mod/:modID", modStatsPage13)
		legacy.GET("/author/", authorPage13)
		legacy.GET("/author/:authorID", authorPage13)
		api := legacy.Group("/api")
		{
			api.GET("/getRandomMod13", getRandomMod13)
			api.GET("/getRandomMod14", getRandomMod14)
		}
	}

	r.GET("/", indexPage14)
	r.GET("/list", modListPage14)
	r.GET("/mod/", modStatsPage14)
	r.GET("/mod/:modID", modStatsPage14)
	r.GET("/author/", authorPage14)
	r.GET("/author/:authorID", authorPage14)
	api := r.Group("/api")
	{
		api.GET("/getRandomMod13", getRandomMod13)
		api.GET("/getRandomMod14", getRandomMod14)
	}

	log.Fatal(r.Run())
}
