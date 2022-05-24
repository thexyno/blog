package main

import (
	"log"
	"net/http"

	"github.com/thexyno/xynoblog"
	"github.com/thexyno/xynoblog/db"
)

func main() {
	port := ":5050"
	log.Printf("Starting Xynoblog on %s", port)
	database := db.NewDb("./blog.db")
	err := database.Seed()
	if err != nil {
		log.Panic(err)
	}
	mux := xynoblog.Mux(database)

	log.Printf("Started Xynoblog on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
