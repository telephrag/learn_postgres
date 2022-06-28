package psql_db

import (
	"fmt"
)

func (pi *PostgresInstance) SampleSelect() error {
	type Row struct {
		Id   int64
		Name string
	}

	rows, err := pi.Conn.Query(
		pi.Ctx,
		"select * from interest order by id",
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	for rows.Next() {
		r := Row{}
		err := rows.Scan(&r.Id, &r.Name)
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
