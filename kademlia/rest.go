package kademlia

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var kademlia *Kademlia

type RequestBody struct {
	Data string `json:"data"`
}

func StartAPI(kad *Kademlia) {
	kademlia = kad
	router := gin.Default()
	router.GET("/objects/:hash", getObject)
	router.POST("/objects", postObject)

	router.Run("0.0.0.0:8080")
}

func getObject(c *gin.Context) {
	hash := c.Param("hash")
	contacts, data := kademlia.LookupData(hash)
	if data != "" {
		c.JSON(http.StatusOK, gin.H{
			"node": contacts[0].String(),
			"data": data,
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Could not find the data",
	})

}

func postObject(c *gin.Context) {
	var requestBody RequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := requestBody.Data
	rsp, err := kademlia.Store([]byte(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"hash": rsp})
}
