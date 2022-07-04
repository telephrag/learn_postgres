package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"postgres_study/changestream"
	"postgres_study/config"
	"postgres_study/expiration"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(), config.PSQL_DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to pgxPool: %v\n", err)
	}
	defer pool.Close()

	csCtx, cancelChangestream := context.WithCancel(context.Background())
	go func() {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Failed to aquire connection from pgxPool: %v\n", err)
		}
		defer conn.Release()

		err = changestream.HandleChangestream(csCtx, conn, "changestream")
		if csCtx.Err() == context.Canceled {
			log.Printf("changestream is closed: %v\n", csCtx.Err())
			return
		}
		if err != nil {
			log.Fatalf("changestream error: %v\n", err)
		}
	}()

	watchCtx, cancelWatch := context.WithCancel(context.Background())
	go func() {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Failed to aquire connection from pgxPool: %v\n", err)
		}
		defer conn.Release()

		err = expiration.WatchExpirations(watchCtx, conn)
		if csCtx.Err() == context.Canceled {
			log.Printf("changestream is closed: %v\n", csCtx.Err())
			return
		}
		if err != nil {
			log.Fatalf("expiration watch error: %v", err)
		}
	}()

	log.Println("Press ^C to cease waiting for notification from changestream and watching expirations.")

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt

	defer fmt.Println()
	defer cancelWatch()
	defer cancelChangestream()

	// wont close right away since WatchExpirations() is waiting
	// for another chance to check if context is done
}
