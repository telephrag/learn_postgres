package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"postgres_study/config"
	"postgres_study/psql_db"
	"syscall"
)

func main() {
	db, err := psql_db.NewInstance(config.PSQL_DB_URL)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v\n", err)
	}
	defer db.Close()

	if err = psql_db.SampleChangestream("changestream"); err != nil {
		log.Fatalf("Changestream error: %v\n", err)
	}

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt

	fmt.Println()
}
