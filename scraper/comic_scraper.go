package scraper

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"komiku-scraper-go/models"
)

func ScrapeComicChapters(db *sql.DB, comicURL string, comicID int64) {
	c := colly.NewCollector(colly.AllowedDomains("komiku.org", "www.komiku.org"))
	c.Limit(&colly.LimitRule{DomainGlob: "*.komiku.*", Parallelism: 1, Delay: 800 * time.Millisecond})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Contains(href, "chapter") {
			title := strings.TrimSpace(e.Text)
			if title == "" {
				title = e.Attr("title")
			}
			if len(title) < 3 {
				return
			}
			slug := strings.Trim(href, "/")
			chapterID := models.EnsureChapter(db, comicID, title, slug, href)
			log.Println("  📖 Chapter:", title)
			ScrapeChapterImages(db, href, chapterID)
		}
	})

	if err := c.Visit(comicURL); err != nil {
		log.Println("Error visiting comic:", err)
	}
}
