package srvone

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/config-loader-concept/pkg/db"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CloseDBConn() {
	r.db.Close()
}

func (r *Repo) Ping() {
	db.Ping(r.db)
}
