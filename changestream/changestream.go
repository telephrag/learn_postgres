package changestream

import (
	"context"
	"encoding/json"
	"log"
	"postgres_study/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

func HandleChangestream(ctx context.Context, conn *pgxpool.Conn, chanName string) error {

	_, err := conn.Exec(ctx, "listen changestream")
	if err != nil {
		return err
	}

	for {
		n, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			return err
		}

		var od models.OplogDiff
		if err = json.Unmarshal([]byte(n.Payload), &od); err != nil {
			return err
		}

		handlerToOp[od.Optype](&od)
	}
}

var handlerToOp = map[string]func(od *models.OplogDiff){
	"INSERT": InsertHandle,
	"UPDATE": nil,
	"DELETE": nil,
}

func InsertHandle(od *models.OplogDiff) {
	log.Println("handling insert...")
	log.Printf("Table name:      %v\n", od.TableName)
	log.Printf("Timestamp:       %v\n", od.Timestamp)
	log.Printf("Inserted record: %v\n", od.New)
}
