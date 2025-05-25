package main

import (
	"log"
	"net"

	"statistics-service/internal/domain/usecases"
	clickhouse "statistics-service/internal/arch/db"
	statsGrpc "statistics-service/internal/interfaces/grpc"
	"statistics-service/proto"

	grpcLib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	clickhouseDSN := "clickhouse:9000"

	repo, err := clickhouse.NewClickHouseStatsRepository(clickhouseDSN)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

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
