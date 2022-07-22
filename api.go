package main

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getRandomMod(c *gin.Context) {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	n := rand.Intn(len(ModList))
	c.String(http.StatusOK, ModList[n].ModName)
}
