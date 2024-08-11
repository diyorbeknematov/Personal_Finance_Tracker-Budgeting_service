package storage

import (
	"budgeting-service/storage/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type IStorage interface {
	AccountRepository() mongodb.AccountRepository
	TransactionRepository() mongodb.TransactionRepository
	BudgetManagementRepo() mongodb.BudgetManagementRepo
	CategoryRepository() mongodb.CategoryRepository
	GoalsRepository() mongodb.GoalsRepository
}

type storageImpl struct {
	// redis *RedisStorage
	mongo *mongo.Database
}

func NewStorage(db *mongo.Database) IStorage {
	return &storageImpl{
		// redis: &RedisStorage{},
		mongo: db,
	}
}

// func (s *storageImpl) GetRedisStorage() *RedisStorage {}

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

func (s *storageImpl) GoalsRepository() mongodb.GoalsRepository {
    return mongodb.NewGoalsRepository(s.mongo)
}