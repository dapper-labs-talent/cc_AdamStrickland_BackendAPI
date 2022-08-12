package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adamstrickland/dapper-api/internal"
	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/adamstrickland/dapper-api/internal/routes"
	"github.com/adamstrickland/dapper-api/internal/users"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xo/dburl"
)

func databaseLocation(cfg *config.Config) string {
	u, err := dburl.Parse(cfg.GetString("databaseUrl"))

	if err != nil {
		log.Fatalf("Bad database URL: %e", err)
	}

	path := u.DSN

	log.Printf("Using database file at '%s'", path)

	return path
}

func Setup(cfg *config.Config) {
	path := databaseLocation(cfg)

	_, err := os.Create(path)

	if err != nil {
		log.Fatalf("Database at '%s' could not be created: %e", path, err)
	}

	log.Printf("Created database at '%s'", path)
}

func Reset(cfg *config.Config) {
	path := databaseLocation(cfg)

	err := os.Remove(path)

	if err != nil {
		log.Fatalf("Database at '%s' could not be removed: %e", path, err)
	}

	log.Printf("Removed database at '%s'", path)

	Setup(cfg)
}

func Migrate(cfg *config.Config) {
	conn, err := internal.NewConnection(cfg)

	if err != nil {
		log.Fatalf("Unable to connect to database: %e", err)
	}

	log.Print("Migrating database:")

	models := []interface{}{
		&users.User{},
	}

	for _, m := range models {
		log.Printf("  Migrating %T", m)
	}

	conn.DB.AutoMigrate(models...)
	log.Print("Migration complete")
}

func Bootstrap(cfg *config.Config, reset bool) {
	if reset {
		Reset(cfg)
	} else {
		Setup(cfg)
	}

	Migrate(cfg)
}

func Run(cfg *config.Config) {
	router := routes.NewRouter(cfg)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.GetString("port")), router))
}

func main() {
	bootstrap := flag.Bool("bootstrap", false, "setup the application")
	reset := flag.Bool("reset", false, "force-recreate the database (if it exists)")
	migrate := flag.Bool("migrate", false, "migrate the database")

	flag.Parse()

	cfg := config.Configuration()

	switch {
	case *bootstrap:
		Bootstrap(cfg, *reset)
	case *migrate:
		Migrate(cfg)
	default:
		Run(cfg)
	}
}
