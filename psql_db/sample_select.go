package psql_db

import (
	"context"
	"fmt"
	"postgres_study/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SampleSelect(ctx context.Context, conn *pgxpool.Conn) error {

	rows, err := conn.Query(
		ctx,
		"select * from interest order by id",
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	for rows.Next() {
		r := models.Interest{}
		err := rows.Scan(r.GetScanForm()...)
		if err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		if rows.Err() != nil {
			return rows.Err()
		}

		fmt.Println(r)
	}

	return nil
}
