package application

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // nolint
	_ "github.com/golang-migrate/migrate/v4/source/file"       // nolint
	"github.com/texazcowboy/warehouse/internal/foundation/database"
	"github.com/texazcowboy/warehouse/internal/foundation/logger"
)

var (
	configPath = os.Getenv("CONFIG")
	sourcePath = os.Getenv("SOURCE")
	direction  = os.Getenv("DIRECTION")
)

func init() {
	flag.Parse()
}

type App struct {
	*Config
	*logger.Logger
	*validator.Validate
	initialized bool
}

func (a *App) Initialize() {
	if a.initialized {
		a.LogEntry.Warn("Application is already initialized")
		return
	}
	a.setupValidator()
	a.readConfiguration()
	a.setupLogger()
}

func (a *App) Run() {
	m, err := setupMigration(a.DBConfig)
	if err != nil {
		a.LogEntry.WithError(err).Fatal("Unable to setup migration")
	}
	err = applyMigration(m)
	if err != nil {
		a.LogEntry.WithError(err).Fatal("Unable to apply migration")
	}
}

func (a *App) readConfiguration() {
	var cfg Config
	if err := cfg.Read(configPath); err != nil {
		panic(err)
	}
	if err := a.Validate.Struct(cfg); err != nil {
		panic(err)
	}
	a.Config = &cfg
}

func (a *App) setupValidator() {
	a.Validate = validator.New()
}

func (a *App) setupLogger() {
	log, err := logger.NewLogger(a.LConfig)
	if err != nil {
		panic(err)
	}
	a.Logger = log
}

func setupMigration(cfg *database.DBConfig) (*migrate.Migrate, error) {
	return migrate.New(sourcePath, buildConnectionString(cfg))
}

func applyMigration(m *migrate.Migrate) error {
	switch direction {
	case "up":
		return m.Up()
	case "down":
		return m.Down()
	default:
		return errors.New("Invalid direction provided: " + direction + " Available directions: [up, down]")
	}
}

func buildConnectionString(cfg *database.DBConfig) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
