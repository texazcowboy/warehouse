package handlers

import (
	"database/sql"

	"github.com/texazcowboy/warehouse/internal/foundation/logger"
)

type Env struct {
	*sql.DB
	*logger.Logger
}

func NewEnvironment(db *sql.DB, logger *logger.Logger) *Env {
	return &Env{DB: db, Logger: logger}
}
