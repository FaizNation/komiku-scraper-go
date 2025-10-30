package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatal("Error connecting to MySQL:", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS komiku")
	if err != nil {
		log.Fatal("Error creating database:", err)
	}

	_, err = db.Exec("USE komiku")
	if err != nil {
		log.Fatal("Error selecting database:", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS types (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS comics (
			id INT AUTO_INCREMENT PRIMARY KEY,
			type_id INT,
			title VARCHAR(255),
			slug VARCHAR(255) UNIQUE,
			url VARCHAR(255),
			author VARCHAR(255),
			description TEXT,
			cover_image VARCHAR(255),
			genres VARCHAR(255),
			status VARCHAR(255),
			release_year VARCHAR(255),
			chapter_count INT,
			FOREIGN KEY (type_id) REFERENCES types(id)
		)`,
		`CREATE TABLE IF NOT EXISTS chapters (
			id INT AUTO_INCREMENT PRIMARY KEY,
			comic_id INT,
			title VARCHAR(255),
			slug VARCHAR(255) UNIQUE,
			url VARCHAR(255),
			chapter_number VARCHAR(255),
			FOREIGN KEY (comic_id) REFERENCES comics(id)
		)`,
		`CREATE TABLE IF NOT EXISTS images (
			id INT AUTO_INCREMENT PRIMARY KEY,
			chapter_id INT,
			position INT,
			url VARCHAR(255) UNIQUE,
			FOREIGN KEY (chapter_id) REFERENCES chapters(id)
		)`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatal("Error creating table:", err)
		}
	}

	log.Println("Database and tables created successfully ✅")
}
