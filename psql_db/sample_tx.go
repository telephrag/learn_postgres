package psql_db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SampleTransactionSuccess(ctx context.Context, conn pgxpool.Conn) error {
	tx, err := conn.Conn().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // does nothing if transaction succeeds, rolls back otherwise

	_, err = tx.Exec(ctx, "insert into interest(name) values('Golang')")
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
