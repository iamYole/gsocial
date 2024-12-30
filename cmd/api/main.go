package main

import (
	"log"

	"github.com/iamYole/gsocial/internal/db"
	"github.com/iamYole/gsocial/internal/env"
	"github.com/iamYole/gsocial/internal/store"
)

const version = "0.0.1"

func main() {
	config := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			dsn:          env.GetString("DB_ADDR", "postgresql://sampleuser:samplepassword@localhost:sampleport/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONN", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONN", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}
	db, err := db.New(config.db.dsn, config.db.maxIdleTime, config.db.maxOpenConns, config.db.maxIdleConns)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Database connected sucessuflly")

	store := store.NewStorage(db)

	app := &application{
		config: config,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
