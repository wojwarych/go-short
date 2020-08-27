package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/warycwoj/go-short/db"
	"github.com/warycwoj/go-short/shortener"
)

type LongURL struct {
	LongURL string `form:"longURL"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	if _, err := db.Open(); err != nil {
		panic(err)
	}
	if !db.DB.HasTable(&URL{}){
		db.DB.AutoMigrate(&URL{})
	}
	defer db.Close()
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
		urls := db.DB.Find(&URL{})
		log.Println(urls)
	})
	r.POST("/", func(c *gin.Context) {
		var longURL LongURL
		c.ShouldBind(&longURL)
		shortURL := shortener.Shorten(12345)
		newURL := &URL{LongURL: longURL.LongURL, ShortURL: shortURL}
		db.DB.Create(newURL)
		log.Printf("%s, %v", "Successfully added new URL to DB", *newURL)
	})
	r.Run(":8000")
}
