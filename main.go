package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/warycwoj/go-short/shortener"
)

type LongURL struct {
	LongURL string `form:"longURL"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	r.POST("/", func(c *gin.Context) {
		var longURL LongURL
		c.ShouldBind(&longURL)
		log.Printf("%v", shortener.Shorten(12345))
	})
	r.Run(":8000")
}
