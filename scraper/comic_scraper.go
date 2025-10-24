package scraper

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"komiku-scraper-go/models"

	"github.com/gocolly/colly/v2"
)

func ScrapeComicDetailsAndChapters(db *sql.DB, comicURL string, comicID int64) {
	c := colly.NewCollector(colly.AllowedDomains("komiku.org", "www.komiku.org"))
	c.Limit(&colly.LimitRule{DomainGlob: "*.komiku.*", Parallelism: 1, Delay: 800 * time.Millisecond})

	var author, description, coverImage, genres, status, releaseYear sql.NullString
	var foundChapterCount int64
	var genreList []string

	c.OnHTML("table.inftable tr", func(e *colly.HTMLElement) {
		label := e.ChildText("td:nth-child(1)")
		value := strings.TrimSpace(e.ChildText("td:nth-child(2)"))

		switch strings.TrimSpace(label) {
		case "Pengarang":
			author = sql.NullString{String: value, Valid: value != ""}
		case "Status":
			status = sql.NullString{String: value, Valid: value != ""}
		case "Rilis":
			releaseYear = sql.NullString{String: value, Valid: value != ""}
		}
	})

	c.OnHTML(".desc", func(e *colly.HTMLElement) {
		descText := strings.TrimSpace(e.Text)
		description = sql.NullString{String: descText, Valid: descText != ""}
	})

	c.OnHTML(".ims img", func(e *colly.HTMLElement) {
		imgSrc := e.Request.AbsoluteURL(e.Attr("src"))
		coverImage = sql.NullString{String: imgSrc, Valid: imgSrc != ""}
	})

	c.OnHTML("ul.genre li.genre", func(e *colly.HTMLElement) {

		genreName := strings.TrimSpace(e.Text)
		if genreName != "" {
			genreList = append(genreList, genreName)
		}
	})

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

			// --- (BARU) Filter 2: Lewati chapter "Awal:" dan "Terbaru:" ---
			if strings.Contains(title, "Awal:") || strings.Contains(title, "Terbaru:") {
				log.Printf("Skip chapter (Awal/Terbaru): %s\n", title)
				return // Lewati chapter ini untuk mencegah duplikasi
			}

			foundChapterCount++

			slug := strings.Trim(href, "/")
			chapterID := models.EnsureChapter(db, comicID, title, slug, href)
			log.Println(" 	📖 Chapter:", title)
			ScrapeChapterImages(db, href, chapterID)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		log.Printf("Selesai scrape detail: %s. Menemukan %d chapter.\n", comicURL, foundChapterCount)

		var chapterCountToUpdate sql.NullInt64
		if foundChapterCount > 0 {
			chapterCountToUpdate = sql.NullInt64{Int64: foundChapterCount, Valid: true}
		}
		if len(genreList) > 0 {
			genreText := strings.Join(genreList, ", ")
			genres = sql.NullString{String: genreText, Valid: true}
		}

		models.UpdateComicDetails(db, comicID, author, description, coverImage, genres, status, releaseYear, chapterCountToUpdate)
	})

	if err := c.Visit(comicURL); err != nil {
		log.Println("Error visiting comic:", err)
	}
}
