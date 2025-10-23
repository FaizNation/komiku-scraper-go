package models

import (
	"database/sql"
	"log"
)

func EnsureType(db *sql.DB, name string) int64 {
	res, err := db.Exec("INSERT IGNORE INTO types (name) VALUES (?)", name)
	if err != nil {
		log.Println("Error insert type:", err)
		return 0
	}
	id, _ := res.LastInsertId()
	if id != 0 {
		return id
	}
	var existing int64
	db.QueryRow("SELECT id FROM types WHERE name = ?", name).Scan(&existing)
	return existing
}
