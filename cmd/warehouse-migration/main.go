package main

import (
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
)

func init() {
	flag.Parse()
}

// note: only goes in the up direction and applies all scripts
func main() {
	var cfg ApplicationConfig
	cfg.Read(configPath)
	m, err := migrate.New(*sourcePath, buildConnectionString(&cfg.DBConfig))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func buildConnectionString(cfg *database.DBConfig) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
