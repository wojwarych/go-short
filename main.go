package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/warycwoj/go-short/db"
	"github.com/warycwoj/go-short/models"
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
	urlTable := models.URL{}
	db.Migrator().DropTable(&urlTable)
	db.Migrator().CreateTable(&urlTable)
	if db.Migrator().HasTable(&urlTable) {
		db.Migrator().AutoMigrate(&urlTable)
	}
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
		rows, _ := db.Model(&urlTable).Rows()
		for rows.Next() {
			var url models.URL
			db.ScanRows(rows, &url)
			log.Printf("%v\n, %d", url, url.ID)
		}
	})
	r.POST("/", func(c *gin.Context) {
		var postedURL PostedURL
		c.ShouldBind(&postedURL)
		lastRecord := models.URL{}
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
		newURLRow := &models.URL{LongURL: postedURL.LongURL, ShortURL: short.String()}
		db.Create(newURLRow)
	})
	r.GET("/:url/", func(c *gin.Context) {
		desiredURL := c.Param("url")
		decodedPK := shortener.Decoder(desiredURL)
		urlModel := models.URL{}
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
