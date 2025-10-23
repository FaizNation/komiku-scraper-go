package models

import (
	"database/sql"
	"log"
)

func EnsureComic(db *sql.DB, typeID int64, title, slug, url string) int64 {
	res, err := db.Exec("INSERT IGNORE INTO comics (type_id, title, slug, url) VALUES (?, ?, ?, ?)",
		typeID, title, slug, url)
	if err != nil {
		log.Println("Error insert comic:", err)
		return 0
	}
	id, _ := res.LastInsertId()
	if id != 0 {
		return id
	}
	var existing int64
	db.QueryRow("SELECT id FROM comics WHERE slug = ?", slug).Scan(&existing)
	return existing
}
