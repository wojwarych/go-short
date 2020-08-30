package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/warycwoj/go-short/db"
	"github.com/warycwoj/go-short/shortener"
)

type PostedURL struct {
	LongURL string `form:"longURL"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	db, err := db.Open()
	if err != nil {
		panic(err)
	}
	db.Migrator().DropTable(&URL{})
	db.Migrator().CreateTable(&URL{})
	if db.Migrator().HasTable(&URL{}) {
		db.Migrator().AutoMigrate(&URL{})
	}
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
		rows, _ := db.Model(&URL{}).Rows()
		for rows.Next() {
			var url URL
			db.ScanRows(rows, &url)
			log.Printf("%v\n, %d", url, url.ID)
		}
	})
	r.POST("/", func(c *gin.Context) {
		var postedURL PostedURL
		c.ShouldBind(&postedURL)
		lastRecord := URL{}
		ret := db.Last(&lastRecord)
		var shortPath string
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			shortPath = shortener.Shorten(1)
		} else {
			shortPath = shortener.Shorten(lastRecord.ID + 1)
		}
		_, err := url.Parse(postedURL.LongURL)
		if err != nil {
			panic(err)
		}
		short := url.URL{
			Scheme: "http",
			Host:   "localhost:8000",
			Path:   shortPath,
		}
		log.Printf("%s", short)
		newURLRow := &URL{LongURL: postedURL.LongURL, ShortURL: short.String()}
		db.Create(newURLRow)
	})
	r.GET("/:url/", func(c *gin.Context) {
		desiredURL := c.Param("url")
		decodedPK := shortener.Decoder(desiredURL)
		urlModel := URL{}
		ret := db.First(&urlModel, decodedPK)
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			c.HTML(http.StatusNotFound, "URL not found in DB!", nil)
		} else {
			c.Redirect(http.StatusFound, urlModel.LongURL)
			db.Model(&urlModel).Update("visits", urlModel.Visits+1)
		}
	})
	r.Run(":8000")
}
