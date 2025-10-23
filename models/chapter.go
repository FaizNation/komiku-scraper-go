package models

import (
	"database/sql"
	"log"
	"strings"
)

func EnsureChapter(db *sql.DB, comicID int64, title, slug, url string) int64 {
	chNum := extractChapterNumber(url)
	res, err := db.Exec(`INSERT IGNORE INTO chapters 
		(comic_id, title, slug, url, chapter_number) VALUES (?, ?, ?, ?, ?)`,
		comicID, title, slug, url, chNum)
	if err != nil {
		log.Println("Error insert chapter:", err)
		return 0
	}

	id, _ := res.LastInsertId()
	if id != 0 {
		return id
	}
	var existing int64
	db.QueryRow("SELECT id FROM chapters WHERE slug = ?", slug).Scan(&existing)
	return existing
}

func extractChapterNumber(url string) string {
	parts := strings.Split(url, "-")
	return parts[len(parts)-1]
}
