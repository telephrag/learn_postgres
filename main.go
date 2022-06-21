package main

import (
	"context"
	"log"
	"postgres_study/config"

	"github.com/jackc/pgx/v4"
)

func main() {
	dbCtx := context.Background()
	db, err := pgx.Connect(dbCtx, config.DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}
	defer db.Close(dbCtx)
}
