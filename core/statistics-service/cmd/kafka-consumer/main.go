package main

import (
	"context"
	"log"
	clickhouse "statistics-service/internal/arch/db"
	"statistics-service/internal/interfaces/kafka"
)

func main() {
	clickhouseDSN := "clickhouse:9000"
	kafkaBrokers := []string{"broker:29092"}
	kafkaTopic := "post_events"
	kafkaGroupID := "stats_group"

	repo, err := clickhouse.NewClickHouseStatsRepository(clickhouseDSN)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	consumer := kafka.NewConsumer(kafkaBrokers, kafkaTopic, kafkaGroupID, repo)

	log.Println("Starting Kafka consumer...")
	if err := consumer.Consume(context.Background()); err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
