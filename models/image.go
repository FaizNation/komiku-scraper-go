package models

import (
	"database/sql"
	"log"
)

func InsertImage(db *sql.DB, chapterID int64, position int, url string) {
	_, err := db.Exec(`INSERT IGNORE INTO images (chapter_id, position, url) VALUES (?, ?, ?)`,
		chapterID, position, url)
	if err != nil {
		log.Println("Error insert image:", err)
	}
}
