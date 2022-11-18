package main

import (
	"net/http"

	"github.com/NotLe0n/ModStats/server/tmlapi13"
	"github.com/gin-gonic/gin"
)

func getRandomMod(c *gin.Context) {
	c.String(http.StatusOK, tmlapi13.GetRandomMod())
}
