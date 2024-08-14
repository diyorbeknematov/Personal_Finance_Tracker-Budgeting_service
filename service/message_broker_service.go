package service

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"budgeting-service/storage"
	"context"
	"encoding/json"
	"log"
	"log/slog"
)

type MsgBrokerService interface {
	CreateTransaction(msg []byte)
	UpdateBudget(msg []byte)
}

type msBorokerServiceImpl struct {
	storage storage.IStorage
	logger  *slog.Logger
}

var ctx = context.Background()

func NewMsgBrokerService(storage storage.IStorage, logger *slog.Logger) MsgBrokerService {
	return &msBorokerServiceImpl{
		storage: storage,
		logger:  logger,
	}
}

func (m *msBorokerServiceImpl) CreateTransaction(msg []byte) {
	log.Println("Requesting to create transaction")
	var transaction pb.CreateTransactionReq
	err := json.Unmarshal(msg, &transaction)
	if err != nil {
		m.logger.Error("Error unmarshalling transaction message", "error", err)
		return
	}
	_, err = m.storage.TransactionRepository().CreateTransaction(ctx, &transaction)
	if err != nil {
		m.logger.Error("Create transaction error", "error", err)
		return
	}
	err = m.storage.AccountBalance().SetBalance(ctx, models.Balance{
		AccountId: transaction.GetAccountId(),
		Balance:   -transaction.GetAmount(),
	})
	if err != nil {
		m.logger.Error("Set balance error", "error", err)
		return
	}
}

func (m *msBorokerServiceImpl) UpdateBudget(msg []byte) {
	log.Println("Requesting to update budget")
	var budget pb.UpdateBudgetReq
	err := json.Unmarshal(msg, &budget)
	if err != nil {
		m.logger.Error("Error unmarshalling budget message", "error", err)
		return
	}
	_, err = m.storage.BudgetManagementRepo().UpdateBudget(ctx, &budget)
	if err != nil {
		m.logger.Error("Update budget error", "error", err)
		return
	}
}
