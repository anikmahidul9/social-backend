package main

import (
	"log"

	"github.com/anikmahidul9/social/internal/auth"
	"github.com/anikmahidul9/social/internal/db"
	"github.com/anikmahidul9/social/internal/env"
	"github.com/anikmahidul9/social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config{
		addr: env.GetString("ADDR", ":8000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "localhost:5432"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		auth: authConfig{
			secret: env.GetString("JWT_SECRET", "my-secret-key"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("Database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: *cfg,
		store:  store,
		jwt: auth.NewJWTAuthenticator(
			cfg.auth.secret,
			"social-api",
		),
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
