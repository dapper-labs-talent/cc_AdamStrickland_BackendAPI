package internal

import (
	"log"

	"github.com/adamstrickland/dapper-api/internal/config"
	"github.com/xo/dburl"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Conn struct {
	*gorm.DB
}

func NewConnection(cfg *config.Config) (*Conn, error) {
	url, err := dburl.Parse(cfg.GetString("databaseUrl"))

	if err != nil {
		log.Printf("Unable to parse database URL: %e", err)
		return nil, err
	}

	gcfg := gorm.Config{}

	db, _ := gorm.Open(sqlite.Open(url.DSN), &gcfg)

	if err != nil {
		log.Printf("Unable to connect ORM to database: %e", err)
		return nil, err
	}

	conn := &Conn{
		DB: db,
	}

	return conn, nil
}
