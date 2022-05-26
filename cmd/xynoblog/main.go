package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/thexyno/xynoblog"
	"github.com/thexyno/xynoblog/db"
)

func main() {
	port := ":5050"
	log.Printf("Starting Xynoblog on %s", port)
	dburi, ok := os.LookupEnv("XYNOBLOG_DATABASE_URI")
	if !ok {
		dburi = "./blog.db"
	}
	database := db.NewDb(dburi)
	err := database.Seed()
	if err != nil {
		log.Panic(err)
	}
	fontdir, ok := os.LookupEnv("XYNOBLOG_FONT_DIR")
	if !ok {
		fontdir = "./fonts"
	}
	cssdir, ok := os.LookupEnv("XYNOBLOG_CSS_DIR")
	if !ok {
		cssdir = "./cssdist"
	}

	log.Printf("Fontdir: %s, CSSDir: %s", fontdir, cssdir)
	mux := xynoblog.Mux(database, fontdir, cssdir)
	log.Printf("Started Xynoblog on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
