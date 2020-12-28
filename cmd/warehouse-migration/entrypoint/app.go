package entrypoint

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // nolint
	_ "github.com/golang-migrate/migrate/v4/source/file"       // nolint
	"github.com/texazcowboy/warehouse/internal/foundation/database"
	"github.com/texazcowboy/warehouse/internal/foundation/logger"
)

var (
	configPath = flag.String("config", "../../config.yaml", "config file location")
	sourcePath = flag.String("src", "file://migrations", "migration source")
	direction  = flag.String("drc", "up", "migration direction")
)

func init() {
	flag.Parse()
}

type App struct {
	ApplicationConfig
	*logger.Logger
	initialized bool
}

func (a *App) Initialize() {
	if a.initialized {
		a.LogEntry.Warn("Application is already initialized")
		return
	}
	a.readConfiguration()
	a.initLogger()
}

func (a *App) Run() {
	m, err := setupMigration(&a.DBConfig)
	if err != nil {
		a.LogEntry.Fatal(err)
	}
	err = applyMigration(m)
	if err != nil {
		a.LogEntry.Fatal(err)
	}
}

func (a *App) readConfiguration() {
	var cfg ApplicationConfig
	cfg.Read(configPath)
	a.ApplicationConfig = cfg
}

func (a *App) initLogger() {
	log, err := logger.NewLogger(&a.LConfig)
	if err != nil {
		panic(err)
	}
	a.Logger = log
}

func setupMigration(cfg *database.DBConfig) (*migrate.Migrate, error) {
	return migrate.New(*sourcePath, buildConnectionString(cfg))
}

func applyMigration(m *migrate.Migrate) error {
	switch *direction {
	case "up":
		return m.Up()
	case "down":
		return m.Down()
	default:
		return errors.New("Invalid direction provided: " + *direction + " Available directions: [up, down]")
	}
}

func buildConnectionString(cfg *database.DBConfig) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
