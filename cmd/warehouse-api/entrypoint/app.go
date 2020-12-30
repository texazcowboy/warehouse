package entrypoint

import (
	"database/sql"
	"flag"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/texazcowboy/warehouse/internal/foundation/logger"

	"github.com/gorilla/mux"
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/handlers"
	"github.com/texazcowboy/warehouse/internal/foundation/database"
)

var (
	configPath = flag.String("config", "../../config.yaml", "config file location")
)

type App struct {
	*ApplicationConfig
	*logger.Logger
	*mux.Router
	*sql.DB
	*validator.Validate
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
	a.setupValidator()
	a.readConfiguration()
	a.setupLogger()
	a.setupRouter()
	a.openDBConnection()
	a.registerHandlers()
	a.initialized = true
}

func (a *App) Run() {
	a.LogEntry.Infof("Starting server on port %v", a.ServerCfg.Port)
	a.LogEntry.Fatal(http.ListenAndServe(`:`+a.ServerCfg.Port, a.Router))
}

func (a *App) setupLogger() {
	log, err := logger.NewLogger(a.LoggerCfg)
	if err != nil {
		panic(err)
	}
	a.Logger = log
}

func (a *App) setupValidator() {
	a.Validate = validator.New()
}

func (a *App) openDBConnection() {
	a.LogEntry.Info("Connecting to database.")
	a.LogEntry.Infof("Connection params [User: %v | Host: %v | Port: %v | Name: %v]",
		a.DatabaseCfg.User, a.DatabaseCfg.Host, a.DatabaseCfg.Port, a.DatabaseCfg.Name)

	db, err := database.OpenConnection(a.DatabaseCfg)
	if err != nil {
		a.LogEntry.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		a.LogEntry.Fatal(err)
	}
	a.DB = db

	a.LogEntry.Info("Successfully connected")
}

func (a *App) setupRouter() {
	a.Router = mux.NewRouter()
}

func (a *App) readConfiguration() {
	var cfg ApplicationConfig
	if err := cfg.Read(configPath); err != nil {
		panic(err)
	}
	if err := a.Validate.Struct(cfg); err != nil {
		panic(err)
	}
	a.ApplicationConfig = &cfg
}

func (a *App) registerHandlers() {
	a.LogEntry.Info("Registering handlers")
	baseHandler := handlers.BaseHandler{DB: a.DB, Logger: a.Logger, Validate: a.Validate}

	itemHandler := handlers.NewItemHandler(&baseHandler)
	a.Router.HandleFunc("/item", itemHandler.CreateItem).Methods("POST")
	a.Router.HandleFunc("/item/{id:[0-9]+}", itemHandler.GetItem).Methods("GET")
	a.Router.HandleFunc("/items", itemHandler.GetItems).Methods("GET")
	a.Router.HandleFunc("/item/{id:[0-9]+}", itemHandler.UpdateItem).Methods("PUT")
	a.Router.HandleFunc("/item/{id:[0-9]+}", itemHandler.DeleteItem).Methods("DELETE")

	userHandler := handlers.NewUserHandler(&baseHandler)
	a.Router.HandleFunc("/user", userHandler.CreateUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")

	a.LogEntry.Info("Handlers successfully registered")
}
