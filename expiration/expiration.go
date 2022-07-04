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

// Realization of TTL index. Deletes records once they expire.
func WatchExpirations(ctx context.Context, conn *pgxpool.Conn) error { // run it async
	for {
		// get row count and init slice for scans
		countRow, err := conn.Conn().Query(ctx, "select count(id) from oplog")
		if err != nil {
			return fmt.Errorf("failed to get row count: %v", err)
		}

		var opCount int
		countRow.Next()
		err = countRow.Scan(&opCount)
		if err != nil {
			return fmt.Errorf("failed to scan row count: %v", err)
		}
		countRow.Close()

		if opCount == 0 {
			continue
		}

		oplog := make([]OplogExpire, opCount)

		// query and scan into slice
		rows, err := conn.Conn().Query(ctx, "select id, expire from oplog")
		if err != nil {
			return fmt.Errorf("failed to get oplog rows: %v", err)
		}

		var i int = 0
		for rows.Next() {
			err := rows.Scan(oplog[i].GetScanForm()...)
			if err != nil {
				return fmt.Errorf("failed to scan row: %v", err)
			}

			if rows.Err() != nil {
				return rows.Err()
			}

			i++
		}
		rows.Close()

		// query deletions for expired records
		now := time.Now().UTC()
		for _, op := range oplog {
			if op.Expire.Before(now) {
				_, err := conn.Exec(ctx, "delete from oplog where id=$1", op.ID)
				if err != nil {
					return fmt.Errorf("failed to delete expired op record: %v", err)
				}
				log.Printf("deleted record at id: %d\n", op.ID)
			}
		}

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
