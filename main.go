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
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&URL{})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
		rows, _ := db.Model(&URL{}).Rows()
		for rows.Next() {
			var url URL
			db.ScanRows(rows, &url)
			log.Printf("%v\n", url)
		}
	})
	r.POST("/", func(c *gin.Context) {
		var longURL LongURL
		c.ShouldBind(&longURL)
		shortURL := shortener.Shorten(12345)
		newURL := &URL{LongURL: longURL.LongURL, ShortURL: shortURL}
		db.Create(newURL)
		log.Printf("%s, %v", "Successfully added new URL to DB", *newURL)
	})
	r.Run(":8000")
}
