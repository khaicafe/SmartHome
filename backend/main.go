package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serve API first (avoid route conflict)
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	// Serve frontend static files
	r.Static("/assets", filepath.Join("..", "frontend", "dist", "assets"))
	r.StaticFile("/vite.svg", filepath.Join("..", "frontend", "dist", "vite.svg"))
	r.StaticFile("/", filepath.Join("..", "frontend", "dist", "index.html"))

	// Fallback to index.html for SPA routes
	r.NoRoute(func(c *gin.Context) {
		uri := c.Request.RequestURI
		if filepath.Ext(uri) == "" {
			c.File(filepath.Join("..", "frontend", "dist", "index.html"))
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	r.Run(":8080")
}
