package servcie6b631a97

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Ping() error {
	err := r.db.Ping()
	if err != nil {
		return fmt.Errorf("ping failed: %v", err)
	}

	log.Printf("Successfully connected and pinged")
	return nil
}
