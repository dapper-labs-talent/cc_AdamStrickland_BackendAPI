package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Root       = filepath.Join(filepath.Dir(b), "../..")
)

type Config struct {
	viper.Viper
}

func Configuration() *Config {
	v := viper.New()

	v.SetDefault("port", "3000")
	v.BindEnv("port", "PORT")
	log.Printf("Listening on port %s", v.GetString("port"))

	v.SetDefault("env", "development")
	v.BindEnv("env", "APP_ENV")
	log.Printf("Starting %s environment", v.GetString("env"))

	v.SetDefault("tokenHeader", "X-Authentication-Token")

	v.SetDefault("secret", "samplesecret")
	v.BindEnv("secret", "APP_SECRET")

	dds := "sqlite"
	ddrp := fmt.Sprintf("./.data/dapper-api_%s.sqlite3", v.GetString("env"))
	ddfp, _ := filepath.Abs(filepath.Join(Root, ddrp))
	ddu := fmt.Sprintf("%s:%s", dds, ddfp)

	v.SetDefault("databaseUrl", ddu)
	v.BindEnv("databaseUrl", "DATABASE_URL")
	log.Printf("Using database URL '%s'", v.GetString("databaseUrl"))

	return &Config{
		Viper: *v,
	}
}
