package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var hotspotStatus = true

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.DisableConsoleColor()
	hotspot := r.Group("/hotspot")
	{
		hotspot.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"status":  hotspotStatus,
			})
		})
		hotspot.POST("/", func(c *gin.Context) {
			hotspotStatus = true
			c.JSON(200, gin.H{
				"success": true,
			})
		})
		hotspot.DELETE("/", func(c *gin.Context) {
			hotspotStatus = false
			c.JSON(200, gin.H{
				"success": true,
			})
		})
	}
	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
