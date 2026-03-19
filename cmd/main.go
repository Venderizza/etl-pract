package main

import (
	"etl-pract/internal/db"
	"etl-pract/internal/etl"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectPostgres()
	mongoDB := db.ConnectMongo()

	intervalSeconds, err := strconv.Atoi(os.Getenv("ETL_INTERVAL_SECONDS"))
	if err != nil || intervalSeconds <= 0 {
		intervalSeconds = 60
	}

	interval := time.Duration(intervalSeconds) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("Starting replication...")

	for {
		select {
		case <-ticker.C:
			log.Println("Running ETL cycle...")
			if err := etl.Replicate(mongoDB); err != nil {
				log.Println("ETL error:", err)
			} else {
				log.Println("ETL cycle completed successfully")
			}
		}
	}
}
