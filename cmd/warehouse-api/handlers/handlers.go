package handlers

import "database/sql"

type Env struct {
	DB *sql.DB
	// tbd logger, etc.
}

func NewEnvironment(db *sql.DB) *Env {
	return &Env{DB: db}
}
