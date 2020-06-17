package main

import (
	"fmt"
	"log"
	"net/http"

	config "recro_demo/config"
	postgres "recro_demo/postgres"
	"recro_demo/website"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

var (
	db        *postgres.DB
	errDbLoad error
	cfg       *config.Config
)

func init() {
	// init function is called by default before main is called
	// purpose to setup connection which are required by server

	cfg = loadConfig()
	ConnectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.PostgrestHost, cfg.PostgrestPort, cfg.PostgresUser,
		cfg.PostgresDatabase, cfg.PostgresPass,
	)
	db, errDbLoad = postgres.InitPostgresDB(ConnectionString)
	// exit if postgres not initialised
	if errDbLoad != nil {
		log.Println(errDbLoad)
		log.Fatal("Postgres not initialised")
	}
	if cfg.Verbose == "true" {
		db.EnableVerboseMode()
	}
}

func loadConfig() *config.Config {
	// loads environment variables into config
	err := godotenv.Load()
	if err != nil {
		log.Println("main.go: No .env file Found")
	}

	cfg := &config.Config{}
	err = env.Parse(cfg)
	// server won't start until a proper configuration is propvided
	if err != nil || len(cfg.PostgrestHost) == 0 {
		log.Println(err)
		log.Println(cfg.PostgrestHost, "No host")
		log.Fatal("Configuration not defined. Check README.md file")
	}
	return cfg
}

func getRoutes(env *config.Env) *chi.Mux {
	r := chi.NewRouter()
	web := &website.Website{
		Env: env,
	}
	r.Route("/", func(r chi.Router) {
		r.Mount("/", web.GetRouter())
	})
	return r
}

func main() {
	env := &config.Env{}
	env.Config = cfg
	env.Db = db
	router := getRoutes(env)
	fmt.Println("hail hydra")
	http.ListenAndServe(":8080", router)
}
