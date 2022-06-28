package psql_db

import (
	"context"
	"log"
	"postgres_study/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SampleChangestream(chanName string) error {
	pool, err := pgxpool.Connect(context.Background(), config.PSQL_DB_URL)
	if err != nil {
		return err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "listen changestream")
	if err != nil {
		return err
	}

	n, err := conn.Conn().WaitForNotification(context.Background())
	if err != nil {
		return err
	}

	log.Printf("Received notification: %v\n", n)

	return nil
}
