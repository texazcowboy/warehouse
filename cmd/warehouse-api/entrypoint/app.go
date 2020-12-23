package entrypoint

import (
	"database/sql"
	"flag"
	"github.com/gorilla/mux"
	"github.com/texazcowboy/warehouse/internal/foundation/database"
	"log"
	"net/http"
)

var configPath = flag.String("config", "../../config.yaml", "config file location")

func init() {
	flag.Parse()
}

type App struct {
	ApplicationConfig
	*mux.Router
	*sql.DB
}

func (a *App) Initialize() {
	var cfg ApplicationConfig
	log.Println("Trying to read application configuration")
	cfg.Read(configPath)
	log.Println("Application configuration was successfully read")

	log.Printf("Connecting to database.\n")
	log.Printf("User: %v | Host: %v | Port: %v | Name: %v\n",
		cfg.DatabaseCfg.User, cfg.DatabaseCfg.Host, cfg.DatabaseCfg.Port, cfg.DatabaseCfg.Name)
	db, err := database.OpenConnection(&cfg.DatabaseCfg)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected")

	a.DB = db
	a.Router = mux.NewRouter()
	a.ApplicationConfig = cfg

	log.Println("Initializing routes")
	a.initRoutes()
}

func (a *App) Run() {
	log.Printf("Starting server on port %v\n", a.ServerCfg.Port)
	log.Fatal(http.ListenAndServe(`:`+a.ServerCfg.Port, a.Router))
}

func (a *App) initRoutes() {
	// tbd
}
