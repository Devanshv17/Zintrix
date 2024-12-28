package utils

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

var DB *sql.DB

// InitializeDatabase initializes the MySQL database connection
func InitializeDatabase(dsn string) {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test the database connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Println("Database connection successfully established")
}
