package main

import (
	"budgeting-service/config"
	"budgeting-service/pkg/logs"
	"budgeting-service/service"
	"budgeting-service/storage"
	"budgeting-service/storage/mongodb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting budgeting-service...")
	logger := logs.InitLogger()
	cfg := config.Load()
	log.Println("Initializing MongoDB connection...")

	db, err := mongodb.ConnectToMongoDB()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	storage := storage.NewStorage(db)

	listener, err := net.Listen("tcp",
		fmt.Sprintf(":%d", cfg.GRPC_PORT),
	)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	grpcServer := grpc.NewServer()

	service := service.NewServiceManager(listener, grpcServer)
	service.RegisterServiceManagerServer(storage, logger)

	logger.Info("Starting gRPC server...")
	logger.Info("Listening on port", "port", cfg.GRPC_PORT)
	service.Start()
}
