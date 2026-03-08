package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/xanderbilla/my-project/internal/db"
	"github.com/xanderbilla/my-project/internal/env"
	"github.com/xanderbilla/my-project/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgres://admin:password@localhost/my-project?sslmode=disable"),
			maxOpenConn: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConn: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "dev"),
	}

	db, err := db.NewDBConn(cfg.db.addr, cfg.db.maxIdleConn, cfg.db.maxOpenConn, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("Database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config:  cfg,
		storage: store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}