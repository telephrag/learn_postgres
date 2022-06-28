package psql_db

import (
	"fmt"
)

func (pi *PostgresInstance) SampleTransactionSuccess() error {
	tx, err := pi.Conn.Begin(pi.Ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(pi.Ctx) // does nothing if transaction succeeds, rolls back otherwise

	_, err = tx.Exec(pi.Ctx, "insert into interest(name) values('Golang')")
	if err != nil {
		return err
	}

	return tx.Commit(pi.Ctx)
}

func (pi *PostgresInstance) SampleTransactionFailure() error {
	tx, err := pi.Conn.Begin(pi.Ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(pi.Ctx) // does nothing if transaction succeeds, rolls back otherwise

	_, err = tx.Exec(pi.Ctx, "insert into interest(name) values('Golang')")

	return fmt.Errorf("abborted: %v", err)
}
