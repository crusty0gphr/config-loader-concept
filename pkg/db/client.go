package db

import (
	"database/sql"
	"log"
)

func Connect(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}

	return db
}

func Ping(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	log.Print("DB successfully connected and pinged")
}
