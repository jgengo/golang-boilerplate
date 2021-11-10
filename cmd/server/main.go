package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/file"
	"github.com/go-ozzo/ozzo-routing/slash"
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
		Addr:    address,
		Handler: buildHandler(cfg),
	}

	go routing.GracefulShutdown(hs, 10*time.Second, log.Printf)
	log.Printf("server %v is running at %v", Version, address)

	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panicln(err)
	}

}

func buildHandler(cfg *config.Config) http.Handler {
	router := routing.New()

	router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		content.TypeNegotiator(content.JSON),
		fault.Recovery(log.Printf),
	)

	api := router.Group("/api")

	api.Get("/users", func(c *routing.Context) error {
		return c.Write("users#index")
	})
	api.Post("/users", func(c *routing.Context) error {
		return c.Write("users#create")
	})
	api.Put(`/users/<id:\d+>`, func(c *routing.Context) error {
		return c.Write("users#update ->" + c.Param("id"))
	})

	// serve index file
	router.Get("/", file.Content("web/ui/index.html"))

	return router

}
