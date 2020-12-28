package entrypoint

import (
	"database/sql"
	"flag"
	"net/http"

	"github.com/texazcowboy/warehouse/internal/foundation/logger"

	"github.com/gorilla/mux"
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/handlers"
	"github.com/texazcowboy/warehouse/internal/foundation/database"
)

var (
	configPath = flag.String("config", "../../config.yaml", "config file location")
)

type App struct {
	ApplicationConfig
	*logger.Logger
	*mux.Router
	*sql.DB
	initialized bool
}

func init() {
	flag.Parse()
}

func (a *App) Initialize() {
	if a.initialized {
		a.LogEntry.Warn("Application is already initialized")
		return
	}
	a.readConfiguration()
	a.initLogger()
	a.initRouter()
	a.openDBConnection()
	a.registerHandlers()
	a.initialized = true
}

func (a *App) Run() {
	a.LogEntry.Infof("Starting server on port %v", a.ServerCfg.Port)
	a.LogEntry.Fatal(http.ListenAndServe(`:`+a.ServerCfg.Port, a.Router))
}

func (a *App) initLogger() {
	log, err := logger.NewLogger(&a.LoggerCfg)
	if err != nil {
		panic(err)
	}
	a.Logger = log
}

func (a *App) openDBConnection() {
	a.LogEntry.Info("Connecting to database.")
	a.LogEntry.Infof("Connection params [User: %v | Host: %v | Port: %v | Name: %v]",
		a.DatabaseCfg.User, a.DatabaseCfg.Host, a.DatabaseCfg.Port, a.DatabaseCfg.Name)

	db, err := database.OpenConnection(&a.DatabaseCfg)
	if err != nil {
		a.LogEntry.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		a.LogEntry.Fatal(err)
	}
	a.DB = db

	a.LogEntry.Info("Successfully connected")
}

func (a *App) initRouter() {
	a.Router = mux.NewRouter()
}

func (a *App) readConfiguration() {
	var cfg ApplicationConfig
	cfg.Read(configPath)
	a.ApplicationConfig = cfg
}

func (a *App) registerHandlers() {
	a.LogEntry.Info("Registering handlers")

	env := handlers.NewEnvironment(a.DB, a.Logger)

	a.Router.HandleFunc("/item", env.CreateItem).Methods("POST")
	a.Router.HandleFunc("/item/{id:[0-9]+}", env.GetItem).Methods("GET")
	a.Router.HandleFunc("/items", env.GetItems).Methods("GET")
	a.Router.HandleFunc("/item/{id:[0-9]+}", env.UpdateItem).Methods("PUT")
	a.Router.HandleFunc("/item/{id:[0-9]+}", env.DeleteItem).Methods("DELETE")

	a.LogEntry.Info("Handlers successfully registered")
}
