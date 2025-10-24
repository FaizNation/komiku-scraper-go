package scraper

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"komiku-scraper-go/models"
)

func ScrapeListByType(db *sql.DB, tipe string) {
	c := colly.NewCollector(
		colly.AllowedDomains("komiku.org", "www.komiku.org"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*.komiku.*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	listURL := fmt.Sprintf("https://komiku.org/daftar-komik/?tipe=%s", url.QueryEscape(tipe))
	typeID := models.EnsureType(db, tipe)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		if strings.Contains(href, "/manga/") || strings.Contains(href, "/manhwa/") || strings.Contains(href, "/manhua/") {
			title := strings.TrimSpace(e.Text)
			if title == "" {
				title = e.Attr("title")
			}
			if title == "" {
				return
			}

			if strings.HasPrefix(href, "/") {
				href = "https://komiku.org" + href
			}

			slug := strings.Trim(strings.TrimPrefix(href, "https://komiku.org/"), "/")

			comicID := models.EnsureComic(db, typeID, title, slug, href)
			log.Printf("📘 Comic: %s (%s)\n", title, href)

			ScrapeComicDetailsAndChapters(db, href, comicID)
		}
	})

	if err := c.Visit(listURL); err != nil {
		log.Println("❌ Error visiting list:", err)
	}
}
