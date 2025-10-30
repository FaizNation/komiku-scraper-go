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

func UpdateComicDetails(db *sql.DB, comicID int64,
	author, description, coverImage, genres, status, releaseYear sql.NullString,
	chapterCount sql.NullInt64) {

	_, err := db.Exec(`
		UPDATE comics 
		SET 
			author = ?, 
			description = ?, 
			cover_image = ?, 
			genres = ?, 
			status = ?, 
			release_year = ?,
			chapter_count = ?
		WHERE id = ?`,
		author, description, coverImage, genres, status, releaseYear, chapterCount, comicID)

	if err != nil {
		log.Printf("❌ Error updating comic details for ID %d: %v\n", comicID, err)
	} else {
		log.Printf("✅ SUCCES!!! | Updated details for comic ID %d\n", comicID)
	}
}
