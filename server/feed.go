package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog/db"
)

func generateFeed(dbc db.DbConn) (*feeds.Feed, error) {
	posts, err := dbc.Posts(-1, 0)
	if err != nil {
		return nil, err
	}
	items := make([]*feeds.Item, len(posts))
	newest := time.Unix(0, 0)
	for i, v := range posts {
		if v.Updated.After(newest) {
			newest = v.Updated
		}
		items[i] = &feeds.Item{
			Id:          string(v.Id),
			Title:       v.Title,
			Link:        &feeds.Link{Href: fmt.Sprint("https://xyno.space/post/", v.Id)},
			Description: string(Render(v)),
			Created:     v.Created,
		}
	}
	return &feeds.Feed{
		Title:       "xyno - Blog",
		Link:        &feeds.Link{Href: "https://xyno.space"},
		Description: "A Blog about Software Engineering, Hardware, NixOS, and more",
		Author:      &feeds.Author{Name: "xyno", Email: "blog@xyno.space"},
		Updated:     newest,
		Copyright:   fmt.Sprint("(c) ", time.Now().Year(), " xyno (Philipp Hochkamp)"),
		Created:     time.Date(2022, 5, 27, 23, 12, 4, 0, time.UTC),
		Items:       items,
	}, nil
}

func renderISE(c *gin.Context, err error) {
	log.WithField("request", c.FullPath()).Error(err)
	c.Status(500)
	fmt.Fprint(c.Writer, "Internal Server Error")
}

func renderRSS(db db.DbConn) func(c *gin.Context) {
	return func(c *gin.Context) {
		feed, err := generateFeed(db)
		if err != nil {
			renderISE(c, err)
			return
		}
		feedrss, err := feed.ToRss()
		if err != nil {
			renderISE(c, err)
			return
		}
		c.Writer.Header().Add("Content-Type", "application/rss+xml")
		c.Writer.WriteString(feedrss)
	}
}

func renderAtom(db db.DbConn) func(c *gin.Context) {
	return func(c *gin.Context) {
		feed, err := generateFeed(db)
		if err != nil {
			renderISE(c, err)
			return
		}
		feedatom, err := feed.ToAtom()
		if err != nil {
			renderISE(c, err)
			return
		}
		c.Writer.Header().Add("Content-Type", "application/atom+xml")
		c.Writer.WriteString(feedatom)
	}
}

func renderJSONFeed(db db.DbConn) func(c *gin.Context) {
	return func(c *gin.Context) {
		feed, err := generateFeed(db)
		if err != nil {
			renderISE(c, err)
			return
		}
		feedjsonfeed, err := feed.ToJSON()
		if err != nil {
			renderISE(c, err)
			return
		}
		c.Writer.Header().Add("Content-Type", "application/feed+json")
		c.Writer.WriteString(feedjsonfeed)
	}
}
