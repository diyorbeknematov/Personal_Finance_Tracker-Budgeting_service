package main

import (
	"budgeting-service/config"
	"budgeting-service/pkg/logs"
	"budgeting-service/queue/kafka/consumer"
	"budgeting-service/service"
	"budgeting-service/storage"
	"budgeting-service/storage/mongodb"
	"budgeting-service/storage/redis"
	"context"
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

	log.Println("Initializing Redis connection...")
	rdb := redis.ConnectToRedis()

	storage := storage.NewStorage(rdb, db)

	listener, err := net.Listen("tcp",
		fmt.Sprintf(":%d", cfg.GRPC_PORT),
	)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	grpcServer := grpc.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		log.Println("Kafka consumer...")
		reader := consumer.NewKafkaConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupId, logger)
		defer reader.Close()
		err := reader.ConsumeMessages(ctx, func(message []byte) {})
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}()

	service := service.NewServiceManager(listener, grpcServer)
	service.RegisterServiceManagerServer(storage, logger)

	logger.Info("Starting gRPC server...")
	logger.Info("Listening on port", "port", cfg.GRPC_PORT)
	service.Start()
}
