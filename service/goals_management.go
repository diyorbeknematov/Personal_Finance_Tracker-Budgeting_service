package service

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/storage"
	"context"
	"log/slog"
)

type GoalsManagementService interface {
	CreateGoal(ctx context.Context, req *pb.CreateGoalReq) (*pb.CreateGoalResp, error)
	UpdateGoal(ctx context.Context, req *pb.UpdateGoalReq) (*pb.UpdateGoalResp, error)
	DeleteGoal(ctx context.Context, req *pb.DeleteGoalReq) (*pb.DeleteGoalResp, error)
	GetGoal(ctx context.Context, req *pb.GetGoalReq) (*pb.GetGoalResp, error)
	GetGoalsList(ctx context.Context, req *pb.GetGoalsReq) (*pb.GetGoalsResp, error)
}

type goalsManagementServiceImpl struct {
	pb.UnimplementedGoalsManagemenServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewGoalsManagementService(storage storage.IStorage, logger *slog.Logger) *goalsManagementServiceImpl {
	return &goalsManagementServiceImpl{
		storage: storage,
		logger:  logger,
	}
}

func (s *goalsManagementServiceImpl) CreateGoal(ctx context.Context, req *pb.CreateGoalReq) (*pb.CreateGoalResp, error) {
	resp, err := s.storage.GoalsRepository().CreateGoal(ctx, req)
	if err != nil {
		s.logger.Error("Create goal error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *goalsManagementServiceImpl) UpdateGoal(ctx context.Context, req *pb.UpdateGoalReq) (*pb.UpdateGoalResp, error) {
	resp, err := s.storage.GoalsRepository().UpdateGoal(ctx, req)
	if err != nil {
		s.logger.Error("Update goal error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *goalsManagementServiceImpl) DeleteGoal(ctx context.Context, req *pb.DeleteGoalReq) (*pb.DeleteGoalResp, error) {
	resp, err := s.storage.GoalsRepository().DeleteGoal(ctx, req)
	if err != nil {
		s.logger.Error("Delete goal error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *goalsManagementServiceImpl) GetGoal(ctx context.Context, req *pb.GetGoalReq) (*pb.GetGoalResp, error) {
	resp, err := s.storage.GoalsRepository().GetGoal(ctx, req)
	if err != nil {
		s.logger.Error("Get goal error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *goalsManagementServiceImpl) GetGoals(ctx context.Context, req *pb.GetGoalsReq) (*pb.GetGoalsResp, error) {
	resp, err := s.storage.GoalsRepository().GetGoalsList(ctx, req)
	if err != nil {
		s.logger.Error("Get goals list error", "error", err)
		return resp, err
	}
	return resp, nil
}
