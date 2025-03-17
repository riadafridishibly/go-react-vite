package main

import (
	"embed"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist
var server embed.FS

//go:embed frontend/dist/index.html
var index []byte

func main() {
	r := gin.Default()

	r.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Use(static.Serve("/", static.EmbedFolder(server, "frontend/dist")))
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.Status(http.StatusOK)
			c.Writer.Write(index)
		}
	})

	r.Run(":8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
