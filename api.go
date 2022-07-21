package main

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

//api handler to get the whole modList
func getModList(c *gin.Context) {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	c.JSON(http.StatusOK, ModList)
}

func getRandomMod(c *gin.Context) {
	dataMutex.Lock()
	defer dataMutex.Unlock()
	n := rand.Intn(len(ModList))
	c.JSON(http.StatusOK, ModList[n].ModName)
}
