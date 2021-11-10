package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/jgengo/golang-boilerplate/internal/config"
)

var Version = "0.0.1"

var flagConfig = flag.String("config", "./config/dev.yml", "path to the config file")
var portConfig = flag.Int("p", 5002, "port server is listening")

func main() {
	flag.Parse()

	cfg, err := config.Load(*flagConfig, *portConfig)
	if err != nil {
		log.Panicf("failed to load application configuration: %s\n", err)
	}

	db, err := dbx.MustOpen("postgres", cfg.DSN)
	if err != nil {
		log.Panicf("failed to connect to the database: %s\n", err)
	}

	db.QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		log.Printf("[%.2fms] Query SQL: %v", float64(t.Milliseconds()), sql)
	}
	db.ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		log.Printf("[%.2fms] Execute SQL: %v", float64(t.Milliseconds()), sql)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr: address,
		// Handler: buildHandler(cfg),
	}

	fmt.Println(hs)
	// ...

}

// func buildHandler(cfg *config.Config) http.Handler {
// 	router := routing.New()

// }
