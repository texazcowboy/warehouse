package application

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/texazcowboy/warehouse/cmd/warehouse-api/common"
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/item"
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/user"

	"github.com/texazcowboy/warehouse/internal/foundation/security"

	"github.com/texazcowboy/warehouse/internal/foundation/web"

	"github.com/go-playground/validator/v10"

	"github.com/texazcowboy/warehouse/internal/foundation/logger"

	"github.com/gorilla/mux"
	"github.com/texazcowboy/warehouse/internal/foundation/database"
)

var configPath = os.Getenv("CONFIG")

type App struct {
	*Config
	*logger.Logger
	*mux.Router
	*sql.DB
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
	var cfg Config
	if err := cfg.Read(configPath); err != nil {
		panic(err)
	}
	if err := a.Validate.Struct(cfg); err != nil {
		panic(err)
	}
	a.Config = &cfg
}

func (a *App) registerHandlers() {
	a.LogEntry.Info("Registering handlers")
	baseHandler := common.BaseHandler{DB: a.DB, Logger: a.Logger, Validate: a.Validate}

	itemHandler := item.NewItemHandler(&baseHandler)
	a.Router.HandleFunc("/item",
		web.AuthenticationMiddleware(itemHandler.CreateItem, a.LogEntry, security.ValidateToken)).Methods("POST")
	a.Router.HandleFunc("/item/{id:[0-9]+}",
		web.AuthenticationMiddleware(itemHandler.GetItem, a.LogEntry, security.ValidateToken)).Methods("GET")
	a.Router.HandleFunc("/items",
		web.AuthenticationMiddleware(itemHandler.GetItems, a.LogEntry, security.ValidateToken)).Methods("GET")
	a.Router.HandleFunc("/item/{id:[0-9]+}",
		web.AuthenticationMiddleware(itemHandler.UpdateItem, a.LogEntry, security.ValidateToken)).Methods("PUT")
	a.Router.HandleFunc("/item/{id:[0-9]+}",
		web.AuthenticationMiddleware(itemHandler.DeleteItem, a.LogEntry, security.ValidateToken)).Methods("DELETE")

	userHandler := user.NewUserHandler(&baseHandler)
	a.Router.HandleFunc("/user", userHandler.CreateUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}",
		web.AuthenticationMiddleware(userHandler.DeleteUser, a.LogEntry, security.ValidateToken)).Methods("DELETE")
	a.Router.HandleFunc("/login", userHandler.Login).Methods("POST")

	a.LogEntry.Info("Handlers successfully registered")
}
