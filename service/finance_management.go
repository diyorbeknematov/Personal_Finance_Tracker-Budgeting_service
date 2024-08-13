package service

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"budgeting-service/storage"
	"context"
	"log/slog"
)

type FinanceManagementService interface {
	CreateAccount(context.Context, *pb.CreateAccountReq) (*pb.CreateAccountResp, error)
	UpdateAccount(context.Context, *pb.UpdateAccountReq) (*pb.UpdateAccountResp, error)
	GetAccount(context.Context, *pb.GetAccountReq) (*pb.GetAccountResp, error)
	GetAccountsList(context.Context, *pb.GetAccountsListReq) (*pb.GetAccountsListResp, error)
	DeleteAccount(context.Context, *pb.DeleteAccountReq) (*pb.DeleteAccountResp, error)
	// Tranzaksiyalarni boshqarish:
	CreateTransaction(context.Context, *pb.CreateTransactionReq) (*pb.CreateTransactionResp, error)
	UpdateTransaction(context.Context, *pb.UpdateTransactionReq) (*pb.UpdateTransactionResp, error)
	GetTransaction(context.Context, *pb.GetTransactionReq) (*pb.GetTransactionResp, error)
	GetTransactionsList(context.Context, *pb.GetTransactionsListReq) (*pb.GetTransactionsListResp, error)
	DeleteTransaction(context.Context, *pb.DeleteTransactionReq) (*pb.DeleteTransactionResp, error)
}

type financeManagementServiceImpl struct {
	pb.UnimplementedFinanceManagementServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewFinanceManagementService(storage storage.IStorage, logger *slog.Logger) *financeManagementServiceImpl {
	return &financeManagementServiceImpl{
		storage: storage,
		logger:  logger,
	}
}

func (s *financeManagementServiceImpl) CreateAccount(ctx context.Context, req *pb.CreateAccountReq) (*pb.CreateAccountResp, error) {
	resp, err := s.storage.AccountRepository().CreateAccount(ctx, req)
	if err != nil {
		s.logger.Error("Create account error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *financeManagementServiceImpl) UpdateAccount(ctx context.Context, req *pb.UpdateAccountReq) (*pb.UpdateAccountResp, error) {
	resp, err := s.storage.AccountRepository().UpdateAccount(ctx, req)
	if err != nil {
		s.logger.Error("Update account error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *financeManagementServiceImpl) DeleteAccount(ctx context.Context, req *pb.DeleteAccountReq) (*pb.DeleteAccountResp, error) {
	resp, err := s.storage.AccountRepository().DeleteAccount(ctx, req)
	if err != nil {
		s.logger.Error("Delete account error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *financeManagementServiceImpl) GetAccount(ctx context.Context, req *pb.GetAccountReq) (*pb.GetAccountResp, error) {
	resp, err := s.storage.AccountRepository().GetAccount(ctx, req)
	if err != nil {
		s.logger.Error("Get account error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *financeManagementServiceImpl) GetAccountsList(ctx context.Context, req *pb.GetAccountsListReq) (*pb.GetAccountsListResp, error) {
	resp, err := s.storage.AccountRepository().GetAccountsList(ctx, req)
	if err != nil {
		s.logger.Error("Get accounts list error", "error", err)
		return resp, err
	}
	return resp, nil
}

// <---------------------------------------------------------------------->

func (s *financeManagementServiceImpl) CreateTransaction(ctx context.Context, req *pb.CreateTransactionReq) (*pb.CreateTransactionResp, error) {
	resp, err := s.storage.TransactionRepository().CreateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Create transaction error", "error", err)
		return resp, err
	}
	err = s.storage.AccountBalance().SetBalance(ctx, models.Balance{
		AccountId: req.GetAccountId(),
		Balance:   -req.GetAmount(),
	})
	if err != nil {
		s.logger.Error("Set balance error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *financeManagementServiceImpl) UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionReq) (*pb.UpdateTransactionResp, error) {
	resp, err := s.storage.TransactionRepository().UpdateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Update transaction error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *financeManagementServiceImpl) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionReq) (*pb.DeleteTransactionResp, error) {
	resp, err := s.storage.TransactionRepository().DeleteTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Delete transaction error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *financeManagementServiceImpl) GetTransaction(ctx context.Context, req *pb.GetTransactionReq) (*pb.GetTransactionResp, error) {
	resp, err := s.storage.TransactionRepository().GetTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Get transaction error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *financeManagementServiceImpl) GetTransactionsList(ctx context.Context, req *pb.GetTransactionsListReq) (*pb.GetTransactionsListResp, error) {
	resp, err := s.storage.TransactionRepository().GetTransactionsList(ctx, req)
	if err != nil {
		s.logger.Error("Get transactions list error", "error", err)
		return resp, err
	}
	return resp, nil
}
