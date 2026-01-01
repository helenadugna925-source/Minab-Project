package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Global DB variable
var DB *sql.DB

// ConnectDB attempts to connect to the database and assigns it to the global DB variable
func ConnectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// 1. Fallback for local development if .env is missing (Optional but recommended)
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "postgres"
	}
	// password usually has no default, but you can set one if you want

	// 2. Build Connection String
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// 3. Open Connection
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// 4. Ping to verify connection is alive
	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("cannot reach database: %w", err)
	}

	// Print the current database name for easier debugging
	var currentDB string
	if err := DB.QueryRow("SELECT current_database()").Scan(&currentDB); err == nil {
		log.Printf("✅ Connected to PostgreSQL successfully — current database: %s", currentDB)
	} else {
		log.Println("✅ Connected to PostgreSQL successfully")
	}
	return DB, nil
}