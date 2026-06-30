package main

import (
	"context"
	"log"
	"os"

	"github.com/EfoJensen/go-rentrospect/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
	}

	databaseUrl := os.Getenv("DATABASE_URL")

	dbPool, err := pgxpool.New(context.Background(), databaseUrl)

	if err != nil {
		log.Fatal("Db err:", err)
	}

	api := api.NewApiServer(":4000", dbPool)
	api.Run()
}