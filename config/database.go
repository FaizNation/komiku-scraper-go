package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/komiku?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB connection error:", err)
	}

	log.Println("Connected to MySQL successfully ✅")
	return db
}
