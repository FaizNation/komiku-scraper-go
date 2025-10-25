package scraper

import (
	"database/sql"
	"log"
	"net/url"
	"strings"
	"time"
	
	"komiku-scraper-go/models"
	"github.com/gocolly/colly/v2"
)

func ScrapeChapterImages(db *sql.DB, chapterURL string, chapterID int64) {
	c := colly.NewCollector(
		colly.AllowedDomains("komiku.org", "www.komiku.org"),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
	c.Limit(&colly.LimitRule{DomainGlob: "*.komiku.*", Parallelism: 1, Delay: 500 * time.Millisecond})

	pos := 0

	c.OnHTML("img", func(e *colly.HTMLElement) {
		img := e.Attr("data-src")
		if img == "" {
			img = e.Attr("src")
		}
		img = e.Request.AbsoluteURL(strings.TrimSpace(img))
		if img == "" {
			return
		}
		if strings.Contains(img, "flagcdn.com") ||
			strings.Contains(img, "/asset/img/komikuplus2.jpg") {

			log.Printf("   🚫 Skipping unwanted image: %s\n", img)
			return // Lewati gambar ini
		}

		if strings.Contains(img, "gambar-id") || strings.HasSuffix(img, ".jpg") || strings.HasSuffix(img, ".png") {
			pos++
			cleaned := cleanURL(img)
			models.InsertImage(db, chapterID, pos, cleaned)
			log.Printf("    🖼️  Image #%d: %s\n", pos, cleaned)
		}
	})

	if err := c.Visit(chapterURL); err != nil {
		log.Println("Error visiting chapter:", err)
	}
}

func cleanURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.RawQuery = ""
	return u.String()
}
