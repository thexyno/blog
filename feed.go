package xynoblog

import (
	"fmt"
	"net/http"
	"time"

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

func renderISE(w http.ResponseWriter, r *http.Request, err error) {
	log.WithField("request", r.URL.Path).Error(err)
	w.WriteHeader(500)
	fmt.Fprint(w, "Internal Server Error")
	return
}

func renderRSS(db db.DbConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := generateFeed(db)
		if err != nil {
			renderISE(w, r, err)
			return
		}
		feedrss, err := feed.ToRss()
		if err != nil {
			renderISE(w, r, err)
			return
		}
		w.Header().Add("Content-Type", "application/rss+xml")
		fmt.Fprint(w, feedrss)
	}
}

func renderAtom(db db.DbConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := generateFeed(db)
		if err != nil {
			renderISE(w, r, err)
			return
		}
		feedatom, err := feed.ToAtom()
		if err != nil {
			renderISE(w, r, err)
			return
		}
		w.Header().Add("Content-Type", "application/atom+xml")
		fmt.Fprint(w, feedatom)
	}
}

func renderJSONFeed(db db.DbConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := generateFeed(db)
		if err != nil {
			renderISE(w, r, err)
			return
		}
		feedjsonfeed, err := feed.ToJSON()
		if err != nil {
			renderISE(w, r, err)
			return
		}
		w.Header().Add("Content-Type", "application/feed+json")
		fmt.Fprint(w, feedjsonfeed)
	}
}
