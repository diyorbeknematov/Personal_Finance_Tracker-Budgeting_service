package storage

import (
	"budgeting-service/storage/mongodb"
	rdb "budgeting-service/storage/redis"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type IStorage interface {
	AccountRepository() mongodb.AccountRepository
	TransactionRepository() mongodb.TransactionRepository
	BudgetManagementRepo() mongodb.BudgetManagementRepo
	CategoryRepository() mongodb.CategoryRepository
	GoalsRepository() mongodb.GoalsRepository
	ReportingRepository() mongodb.ReportingRepository
	NotificationRepository() mongodb.NotificationRepository
	AccountBalance() rdb.AccountBalanceRepository
}

type storageImpl struct {
	redis *redis.Client
	mongo *mongo.Database
}

func NewStorage(client *redis.Client, db *mongo.Database) IStorage {
	return &storageImpl{
		redis: client,
		mongo: db,
	}
}

func (s *storageImpl) AccountBalance() rdb.AccountBalanceRepository {
	return rdb.NewAccountBalance(s.redis)
}

func (s *storageImpl) AccountRepository() mongodb.AccountRepository {
	return mongodb.NewAccountRepository(s.mongo)
}

func (s *storageImpl) TransactionRepository() mongodb.TransactionRepository {
	return mongodb.NewTransactionRepository(s.mongo)
}

func (s *storageImpl) BudgetManagementRepo() mongodb.BudgetManagementRepo {
	return mongodb.NewBudgetManagementRepo(s.mongo)
}

func (s *storageImpl) CategoryRepository() mongodb.CategoryRepository {
	return mongodb.NewCategoryRepository(s.mongo)
}

func (s *storageImpl) ReportingRepository() mongodb.ReportingRepository {
	return mongodb.NewReportingRepository(s.mongo)
}

func (s *storageImpl) NotificationRepository() mongodb.NotificationRepository {
	return mongodb.NewNotificationRepository(s.mongo)
}

func (s *storageImpl) GoalsRepository() mongodb.GoalsRepository {
	return mongodb.NewGoalsRepository(s.mongo)
}
