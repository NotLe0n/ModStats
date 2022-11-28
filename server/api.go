package main

import (
	"github.com/NotLe0n/ModStats/server/tmlapi14"
	"net/http"

	"github.com/NotLe0n/ModStats/server/tmlapi13"
	"github.com/gin-gonic/gin"
)

func getRandomMod13(c *gin.Context) {
	c.String(http.StatusOK, tmlapi13.GetRandomMod())
}

func getRandomMod14(c *gin.Context) {
	c.String(http.StatusOK, tmlapi14.GetRandomMod())
}
