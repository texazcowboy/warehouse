package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/texazcowboy/warehouse/internal/foundation/database"
	"log"
)

var (
	configPath = flag.String("config", "../../config.yaml", "config file location")
	sourcePath = flag.String("src", "file://migrations", "migration source")
	direction  = flag.String("drc", "up", "migration direction")
)

func init() {
	flag.Parse()
}

// note: applies all scripts
func main() {
	m, err := setupMigration()
	if err != nil {
		log.Fatal(err)
	}
	err = applyMigration(m)
	if err != nil {
		log.Fatal(err)
	}
}

func setupMigration() (*migrate.Migrate, error) {
	var cfg ApplicationConfig
	cfg.Read(configPath)
	return migrate.New(*sourcePath, buildConnectionString(&cfg.DBConfig))
}

func applyMigration(m *migrate.Migrate) error {
	switch *direction {
	case "up":
		return m.Up()
	case "down":
		return m.Down()
	default:
		return errors.New("Invalid direction provided: " + *direction + " Available directions: up, down")
	}
}

func buildConnectionString(cfg *database.DBConfig) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
