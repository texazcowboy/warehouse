package database

import (
	"database/sql"
	"fmt"

	// register drivers.
	_ "github.com/lib/pq"
)

func OpenConnection(cfg *DBConfig) (*sql.DB, error) {
	connInfo := buildConnectionInformation(cfg)
	return sql.Open("postgres", connInfo)
}

func buildConnectionInformation(c *DBConfig) string {
	var format = "user=%s password=%s dbname=%s host=%s port=%s sslmode=disable"
	return fmt.Sprintf(format, c.User, c.Password, c.Name, c.Host, c.Port)
}
