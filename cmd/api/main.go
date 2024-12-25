package main

import (
	"log"

	"github.com/iamYole/gsocial/internal/env"
)

func main() {
	config := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	app := &application{
		config: config,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
