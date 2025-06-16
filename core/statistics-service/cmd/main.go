package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"statistics-service/internal/domain/usecases"
	clickhouse "statistics-service/internal/arch/db"
	statsGrpc "statistics-service/internal/interfaces/grpc"
	"statistics-service/internal/interfaces/kafka"
	"statistics-service/proto"

	grpcLib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		startGRPCServer(repo)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startKafkaConsumer(ctx, kafkaBrokers, kafkaTopic, kafkaGroupID, repo)
	}()

	<-sigChan
	log.Println("Received termination signal, shutting down...")
	cancel()

	wg.Wait()
	log.Println("Service stopped gracefully")
}

func startGRPCServer(repo *clickhouse.ClickHouseStatsRepository) {
	uc := usecases.NewStatsUseCase(repo)
	grpcServer := statsGrpc.NewStatsServer(uc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpcLib.NewServer()
	proto.RegisterStatsServiceServer(s, grpcServer)
	reflection.Register(s)

	log.Println("gRPC server started on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startKafkaConsumer(ctx context.Context, brokers []string, topic, groupID string, repo *clickhouse.ClickHouseStatsRepository) {
	consumer := kafka.NewConsumer(brokers, topic, groupID, repo)

	log.Println("Starting Kafka consumer...")
	if err := consumer.Consume(ctx); err != nil {
		log.Printf("Consumer error: %v", err)
	}
}
