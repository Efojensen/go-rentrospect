package main

import (
	"context"
	"fmt"
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
		log.Fatal("Db err: ", err)
	}

	defer dbPool.Close()

	var greeting string
	err = dbPool.QueryRow(context.Background(), "select 'pg database status: ✅'").Scan(&greeting)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	} else {
		log.Println("database status: ✅")
	}

	api := api.NewApiServer(":4000", dbPool)
	api.Run()
}
