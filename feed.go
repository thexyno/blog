package xynoblog

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog/db"
)

func generateFeed(dbc db.DbConn) (*feeds.Feed, error) {
	sposts, err := dbc.ShortPosts(200, -1, 0)
	if err != nil {
		return nil, err
	}
	items := make([]*feeds.Item, len(sposts))
	newest := time.Unix(0, 0)
	for i, v := range sposts {
		if v.Updated.After(newest) {
			newest = v.Updated
		}
		items[i] = &feeds.Item{
			Id:          string(v.Id),
			Title:       v.Title,
			Link:        &feeds.Link{Href: fmt.Sprint("https://xyno.space/post/", v.Id)},
			Description: string(Render([]byte(v.Content))),
			Created:     v.Created,
		}
	}
	return &feeds.Feed{
		Title:       "xyno - Blog",
		Link:        &feeds.Link{Href: "https://xyno.space"},
		Description: "posts about tech n stuff",
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
	return
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
		fmt.Fprint(c.Writer, feedrss)
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
		fmt.Fprint(c.Writer, feedatom)
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
		fmt.Fprint(c.Writer, feedjsonfeed)
	}
}
