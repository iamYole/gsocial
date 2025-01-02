package main

import (
	"log"

	"github.com/iamYole/gsocial/internal/db"
	"github.com/iamYole/gsocial/internal/env"
	"github.com/iamYole/gsocial/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgresql://yole:sosodef@localhost:5433/social?sslmode=disable")
	conn, err := db.New(addr, "15m", 3, 3)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store, conn)
}
