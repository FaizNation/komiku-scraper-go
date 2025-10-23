package main

import (
	"log"

	"komiku-scraper-go/config"
	"komiku-scraper-go/scraper"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	types := []string{"manga", "manhwa", "manhua"}

	for _, t := range types {
		log.Println("=== Scraping:", t, "===")
		scraper.ScrapeListByType(db, t)
	}
}
