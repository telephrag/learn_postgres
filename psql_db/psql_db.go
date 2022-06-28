package psql_db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type PostgresInstance struct {
	Conn *pgx.Conn
	Ctx  context.Context
}

func NewInstance(url string) (*PostgresInstance, error) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	return &PostgresInstance{
		Conn: db,
		Ctx:  ctx,
	}, nil
}

func (pi *PostgresInstance) Close() {
	pi.Conn.Close(
		pi.Ctx,
	)
}
