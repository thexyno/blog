package server

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog/db"
)

type urlset struct {
	Url   []urlblock `xml:"url"`
	Xmlns string     `xml:"xmlns,attr"`
}

type urlblock struct {
	Loc        string  `xml:"loc"`
	Lastmod    string  `xml:"lastmod"`
	Changefreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

func renderSitemap(db db.DbConn) func(*gin.Context) {
	return func(c *gin.Context) {
		post, err := db.PostIds()
		if err != nil {
			log.WithField("request", c.Request.URL.Path).Error(err)
			renderISE(c, err)
			return
		}
		var newest time.Time
		for _, v := range post {
			if newest.Before(v.UpdatedAt) {
				newest = v.UpdatedAt
			}
		}

		blocks := []urlblock{
			{"https://xyno.space/", time.Now().Format("2006-01-02"), "monthly", 1.0},
			{"https://xyno.space/posts", newest.Format("2006-01-02"), "monthly", 0.2},
		}

		for i := range post {
			time := post[i].UpdatedAt
			id := post[i].PostID
			blocks = append(blocks, urlblock{fmt.Sprintf("https://xyno.space/post/%s", id), time.Format("2006-01-02"), "yearly", 0.9})
		}

		urls := urlset{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", Url: blocks}
		bytes, err := xml.Marshal(urls)
		if err != nil {
			renderISE(c, err)
			return
		}
		c.Writer.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
		c.Writer.Write(bytes)
	}
}
