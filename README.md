# komiku-scraper-go

Go scraper for collecting comic, chapter, and image data from [komiku.org](https://komiku.org) into MySQL.

## Features

- Scrapes comic list pages by type
- Stores comic metadata (title, author, status, genres, release year, cover, description)
- Scrapes chapters for each comic
- Scrapes and stores image URLs for each chapter
- Uses `INSERT IGNORE` to avoid duplicate records

## Project Structure

- `/main.go` - main scraper entry point
- `/database/main.go` - creates database and required tables
- `/config/database.go` - MySQL connection helper
- `/scraper` - scraping logic (list, comic details, chapter images)
- `/models` - DB insert/update helpers

## Requirements

- Go `1.25.1` (see `go.mod`)
- MySQL running on `127.0.0.1:3306`
- MySQL user `root` with empty password (default DSN in code)

## Setup

1. Clone the repository.
2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Create database and tables:

   ```bash
   go run /home/runner/work/komiku-scraper-go/komiku-scraper-go/database/main.go
   ```

4. Run the scraper:

   ```bash
   go run /home/runner/work/komiku-scraper-go/komiku-scraper-go/main.go
   ```

## Database Schema

Tables created by `database/main.go`:

- `types`
- `comics`
- `chapters`
- `images`

## Notes

- Scraping depends on the current HTML structure of `komiku.org`; selector changes may require code updates.
- The database DSN is currently hardcoded in `config/database.go`.
