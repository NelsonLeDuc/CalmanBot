package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	openDB()
}

func openDB() {
	dbURL := os.Getenv("DATABASE_URL")
	database, err := sql.Open("postgres", dbURL)
	if Configuration().VerboseMode() {
		fmt.Println("=== Setting up DB ===")
		fmt.Println("    DB: ", database)
		fmt.Printf("DB err: %v\n\n", err)
	}
	if err != nil {
		log.Fatalf("[x] Could not open the connection to the database. Reason: %s", err.Error())
	}

	db = database
}

func DB() *sql.DB {
	stats := db.Stats()
	if stats.OpenConnections >= 20 {
		db.Close()
		log.Println("Database has too many connections, closing and reopening")
		openDB()
	}

	return db
}
