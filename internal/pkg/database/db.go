package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func New() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "data/main.db")
	if err != nil {
		panic(err)
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	return db
}
