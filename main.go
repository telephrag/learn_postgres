package main

import (
	"context"
	"errors"
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
		log.Panicf("Failed to connect to pgxPool: %v\n", err)
	}
	defer pool.Close()

	csCtx, cancelChangestream := context.WithCancel(context.Background())
	go func() {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Panicf("Failed to aquire connection from pgxPool: %v\n", err)
		}
		defer conn.Release()

		err = changestream.HandleChangestream(csCtx, conn, "changestream")
		if err != nil && errors.Unwrap(err) != context.Canceled { // do not lose context's error
			defer log.Panicf("changestream error: %v\n", err)
		}
		if csCtx.Err() == context.Canceled {
			log.Printf("changestream is closed: %v\n", csCtx.Err())
			return
		}
	}()

	watchCtx, cancelWatch := context.WithCancel(context.Background())
	go func() {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Panicf("Failed to aquire connection from pgxPool: %v\n", err)
		}
		defer conn.Release()

		err = expiration.WatchExpirations(watchCtx, conn)

		if err != nil && errors.Unwrap(err) != context.Canceled {
			defer log.Panicf("expiration watch error: %v", err)
		}
		if csCtx.Err() == context.Canceled {
			log.Printf("deletion of expired records is ceased: %v\n", csCtx.Err())
			return
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
