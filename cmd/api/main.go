package main

import (
	"time"

	"github.com/iamYole/gsocial/internal/db"
	"github.com/iamYole/gsocial/internal/env"
	"github.com/iamYole/gsocial/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			gSocial Backend
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	config := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			dsn:          env.GetString("DB_ADDR", "postgresql://sampleuser:samplepassword@localhost:sampleport/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONN", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONN", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour * 24,
		},
	}

	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database Connection
	db, err := db.New(config.db.dsn, config.db.maxIdleTime, config.db.maxOpenConns, config.db.maxIdleConns)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer db.Close()
	logger.Info("Database connected sucessuflly")

	store := store.NewStorage(db)

	app := &application{
		config: config,
		store:  store,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
