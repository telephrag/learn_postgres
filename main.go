package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"postgres_study/config"
	"postgres_study/psql_db"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	pool, err := pgxpool.Connect(ctx, config.PSQL_DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to pgxPool: %v\n", err)
	}
	defer pool.Close()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("Failed to aquire connection from pgxPool: %v\n", err)
	}
	defer conn.Conn().Close(ctx)

	if err = psql_db.SampleSelect(ctx, conn); err != nil {
		log.Fatalf("Failed select query: %v\n", err)
	}

	go func() { // we are not using this connection anywhere else so async is fine
		if err = psql_db.SampleChangestream(ctx, conn, "changestream"); err != nil {
			if ctx.Err() == context.Canceled {
				log.Printf("Changestream is closed: %v\n", ctx.Err())
			} else {
				log.Fatalf("Changestream error: %v\n", err)
			}
		}
	}()

	log.Println("Press ^C to ceasy waiting for notification from changestream.")

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt
	cancel()

	<-interupt

	fmt.Println()
	os.Exit(0)
}
