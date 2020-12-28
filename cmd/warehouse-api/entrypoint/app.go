package entrypoint

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/handlers"
	"github.com/texazcowboy/warehouse/internal/foundation/database"
)

var (
	configPath = flag.String("config", "../../config.yaml", "config file location")
)

type App struct {
	ApplicationConfig
	*mux.Router
	*sql.DB
	initialized bool
}

func init() {
	flag.Parse()
}

func (a *App) Initialize() {
	if a.initialized {
		log.Println("Application is already initialized")
		return
	}
	a.readConfiguration()
	a.openDBConnection()
	a.initRouter()
	a.registerHandlers()
	a.initialized = true
}

func (a *App) Run() {
	log.Printf("Starting server on port %v\n", a.ServerCfg.Port)
	log.Fatal(http.ListenAndServe(`:`+a.ServerCfg.Port, a.Router))
}

func (a *App) openDBConnection() {
	log.Printf("Connecting to database.\n")
	log.Printf("Connection params [User: %v | Host: %v | Port: %v | Name: %v]\n",
		a.DatabaseCfg.User, a.DatabaseCfg.Host, a.DatabaseCfg.Port, a.DatabaseCfg.Name)

	db, err := database.OpenConnection(&a.DatabaseCfg)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	a.DB = db

	log.Println("Successfully connected")
}

func (a *App) initRouter() {
	a.Router = mux.NewRouter()
}

func (a *App) readConfiguration() {
	log.Println("Trying to read application configuration")

	var cfg ApplicationConfig
	cfg.Read(configPath)
	a.ApplicationConfig = cfg

	log.Println("Application configuration successfully read")
}

func (a *App) registerHandlers() {
	log.Println("Registering handlers")

	env := handlers.NewEnvironment(a.DB)

	a.Router.HandleFunc("/item", env.CreateItem).Methods("POST")
	a.Router.HandleFunc("/item/{id:[0-9]+}", env.GetItem).Methods("GET")
	a.Router.HandleFunc("/items", env.GetItems).Methods("GET")
	a.Router.HandleFunc("/item/{id:[0-9]+}", env.UpdateItem).Methods("PUT")
	a.Router.HandleFunc("/item/{id:[0-9]+}", env.DeleteItem).Methods("DELETE")

	log.Println("Handlers successfully registered")
}
