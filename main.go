package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rylio/ytdl"
)

func audioById(id string) (*url.URL, error, string) {
	vid, err := ytdl.GetVideoInfo(id)
	if err != nil {
		fmt.Println(err)
	}
	filtered := vid.Formats.Filter(ytdl.FormatItagKey, []interface{}{"139", "140", "141", "171", "172", "249", "250", "251"})
	rnd := filtered[rand.Intn(len(filtered))]
	url, err := vid.GetDownloadURL(rnd)
	return url, err, vid.Title
}

func handleWebError(err error, c *gin.Context) {
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(400, err)
	}
}

type track struct {
	Url   string `json:"id"`
	Title string `json:"title"`
}

func web() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/video", func(c *gin.Context) {
		id := c.DefaultQuery("id", "CJyXfq2xdYY")
		fmt.Println("get video id=")
		url, err, _ := audioById(id)
		handleWebError(err, c)
		c.Redirect(http.StatusMovedPermanently, url.String())
	})

	router.GET("/playlist", func(c *gin.Context) {
		plid := c.DefaultQuery("id", "PLVWQmUmEmGaAaq3JoDGw02_whc-eCWJY0")
		ids := idsFromPlaylist(plid)
		var tracks []track

		for _, i := range ids {
			url, err, title := audioById(i)
			handleWebError(err, c)
			tracks = append(tracks, track{Url: url.String(), Title: title})
		}

		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"tracks": tracks,
		})
	})

	router.Run(":" + port)
}

func main() {
	web()
}
