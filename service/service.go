package service

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/storage"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type ServiceManager interface {
	RegisterServiceManagerServer(storage storage.IStorage, logger *slog.Logger)
	Start() error
}

type serviceManagerImpl struct {
	listener net.Listener
	server   *grpc.Server
}

func NewServiceManager(listener net.Listener, server *grpc.Server) ServiceManager {
	return &serviceManagerImpl{
		listener: listener,
		server:   server,
	}
}

func (sm *serviceManagerImpl) RegisterServiceManagerServer(storage storage.IStorage, logger *slog.Logger) {
	log.Println("Registering budgeting-service")

	pb.RegisterBudgetingServiceServer(sm.server, NewBudgetManagementService(storage, logger))
	pb.RegisterFinanceManagementServiceServer(sm.server, NewFinanceManagementService(storage, logger))
	pb.RegisterGoalsManagemenServiceServer(sm.server, NewGoalsManagementService(storage, logger))
	pb.RegisterReportingNotificationServiceServer(sm.server, NewReportingNotificationService(storage, logger))
}

func (sm *serviceManagerImpl) Start() error {
	log.Println("Running budgeting-service")
	log.Printf("Server is running %v", sm.listener.Addr())

	return sm.server.Serve(sm.listener)
}
