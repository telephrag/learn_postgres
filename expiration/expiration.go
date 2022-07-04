package expiration

import (
	"context"
	"fmt"
	"log"
	"postgres_study/config"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type OplogExpire struct {
	ID     int64
	Expire time.Time
}

func (o *OplogExpire) GetScanForm() []interface{} {
	return []interface{}{
		&o.ID,
		&o.Expire,
	}
}

func WatchExpirations(ctx context.Context, conn *pgxpool.Conn) error {
	for {
		now := time.Now().UTC()
		ct, err := conn.Exec(context.Background(), "delete from oplog where expire<$1", now)
		if err != nil {
			return fmt.Errorf("failed to delete expired op records: %v", err)
		}
		log.Printf("%d expired records deleted from oplog\n", ct.RowsAffected())

		sleep := time.After(time.Second * config.WATCH_EXPIRATIONS_INTERVAL_SECONDS)
	hang:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err() // context.Cancelled
			case <-sleep:
				break hang
			}
		}
	}
}
