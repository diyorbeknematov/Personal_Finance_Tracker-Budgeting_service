package service

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/storage"
	"context"
	"log/slog"
)

type BudgetManagementService interface {
	CreateCategory(context.Context, *pb.CreateCategoryReq) (*pb.CreateCategoryResp, error)
	GetCategoriesList(context.Context, *pb.GetCategoriesReq) (*pb.GetCategoriesResp, error)
	GetCategory(context.Context, *pb.GetCategoryReq) (*pb.GetCategoryResp, error)
	UpdateCategory(context.Context, *pb.UpdateCategoryReq) (*pb.UpdateCategoryResp, error)
	DeleteCategory(context.Context, *pb.DeleteCategoryReq) (*pb.DeleteCategoryResp, error)
	// Byudjetni boshqarish
	CreateBudget(context.Context, *pb.CreateBudgetReq) (*pb.CreateBudgetResp, error)
	GetBudgetsList(context.Context, *pb.GetBudgetsReq) (*pb.GetBudgetsResp, error)
	GetBudget(context.Context, *pb.GetBudgetReq) (*pb.GetBudgetResp, error)
	UpdateBudget(context.Context, *pb.UpdateBudgetReq) (*pb.UpdateBudgetResp, error)
	DeleteBudget(context.Context, *pb.DeleteBudgetReq) (*pb.DeleteBudgetResp, error)
}

type budgetManagementServiceImpl struct {
	pb.UnimplementedBudgetingServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewBudgetManagementService(storage storage.IStorage, logger *slog.Logger) *budgetManagementServiceImpl {
	return &budgetManagementServiceImpl{
		storage: storage,
		logger:  logger,
	}
}

func (s *budgetManagementServiceImpl) CreateBudget(ctx context.Context, req *pb.CreateBudgetReq) (*pb.CreateBudgetResp, error) {
	resp, err := s.storage.BudgetManagementRepo().CreateBudget(ctx, req)
	if err != nil {
		s.logger.Error("Create budget error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *budgetManagementServiceImpl) UpdateBudget(ctx context.Context, req *pb.UpdateBudgetReq) (*pb.UpdateBudgetResp, error) {
	resp, err := s.storage.BudgetManagementRepo().UpdateBudget(ctx, req)
	if err != nil {
		s.logger.Error("Update budget error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *budgetManagementServiceImpl) DeleteBudget(ctx context.Context, req *pb.DeleteBudgetReq) (*pb.DeleteBudgetResp, error) {
	resp, err := s.storage.BudgetManagementRepo().DeleteBudget(ctx, req)
	if err != nil {
		s.logger.Error("Delete budget error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *budgetManagementServiceImpl) GetBudget(ctx context.Context, req *pb.GetBudgetReq) (*pb.GetBudgetResp, error) {
	resp, err := s.storage.BudgetManagementRepo().GetBudget(ctx, req)
	if err != nil {
		s.logger.Error("Get budget error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *budgetManagementServiceImpl) GetBudgetsList(ctx context.Context, req *pb.GetBudgetsReq) (*pb.GetBudgetsResp, error) {
	resp, err := s.storage.BudgetManagementRepo().GetBudgetsList(ctx, req)
	if err != nil {
		s.logger.Error("Get budgets list error", "error", err)
		return resp, err
	}
	return resp, nil
}

// <------------------------------------------------------------------------>

func (s *budgetManagementServiceImpl) CreateCategory(ctx context.Context, req *pb.CreateCategoryReq) (*pb.CreateCategoryResp, error) {
	resp, err := s.storage.CategoryRepository().CreateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Create category error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *budgetManagementServiceImpl) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryReq) (*pb.UpdateCategoryResp, error) {
	resp, err := s.storage.CategoryRepository().UpdateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Update category error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *budgetManagementServiceImpl) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryReq) (*pb.DeleteCategoryResp, error) {
	resp, err := s.storage.CategoryRepository().DeleteCategory(ctx, req)
	if err != nil {
		s.logger.Error("Delete category error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *budgetManagementServiceImpl) GetCategory(ctx context.Context, req *pb.GetCategoryReq) (*pb.GetCategoryResp, error) {
	resp, err := s.storage.CategoryRepository().GetCategory(ctx, req)
	if err != nil {
		s.logger.Error("Get category error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *budgetManagementServiceImpl) GetCategoriesList(ctx context.Context, req *pb.GetCategoriesReq) (*pb.GetCategoriesResp, error) {
	resp, err := s.storage.CategoryRepository().GetCategoriesList(ctx, req)
	if err != nil {
		s.logger.Error("Get categories list error", "error", err)
		return resp, err
	}
	return resp, nil
}
