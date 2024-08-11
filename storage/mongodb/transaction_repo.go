package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"budgeting-service/pkg/enums"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *pb.CreateTransactionReq) (*pb.CreateTransactionResp, error)
	UpdateTransaction(ctx context.Context, transaction *pb.UpdateTransactionReq) (*pb.UpdateTransactionResp, error)
	DeleteTransaction(ctx context.Context, request *pb.DeleteTransactionReq) (*pb.DeleteTransactionResp, error)
	GetTransaction(ctx context.Context, request *pb.GetTransactionReq) (*pb.GetTransactionResp, error)
	GetTransactionsList(ctx context.Context, request *pb.GetTransactionsListReq) (*pb.GetTransactionsListResp, error)
}

type transactionRepositoryImpl struct {
	coll *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database) TransactionRepository {
	return &transactionRepositoryImpl{coll: db.Collection("transactions")}
}

func (repo *transactionRepositoryImpl) CreateTransaction(ctx context.Context, transaction *pb.CreateTransactionReq) (*pb.CreateTransactionResp, error) {
	date, err := time.Parse("2006-01-02 15:04:05", transaction.Date)
	if err != nil {
		return nil, err
	}
	_, err = repo.coll.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.NewString()},
		{Key: "account_id", Value: transaction.AccountId},
		{Key: "user_id", Value: transaction.UserId},
		{Key: "type", Value: transaction.Type.String()},
		{Key: "amount", Value: transaction.Amount},
		{Key: "description", Value: transaction.Description},
		{Key: "date", Value: date},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
		{Key: "deleted_at", Value: nil},
	})

	if err != nil {
		return &pb.CreateTransactionResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}
	return &pb.CreateTransactionResp{
		Status:  "success",
		Message: "created transaction successfully",
	}, nil
}

func (repo *transactionRepositoryImpl) UpdateTransaction(ctx context.Context, transaction *pb.UpdateTransactionReq) (*pb.UpdateTransactionResp, error) {
	updateDate, err := time.Parse("2006-01-02 15:04:05", transaction.Date)
	if err != nil {
		return nil, err
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "type", Value: transaction.Type},
			{Key: "amount", Value: transaction.Amount},
			{Key: "description", Value: transaction.Description},
			{Key: "date", Value: updateDate},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	filter := bson.D{
		{Key: "_id", Value: transaction.Id},
		{Key: "deleted_at", Value: nil},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return &pb.UpdateTransactionResp{
			Status:  "error",
			Message: "Error updating transaction: " + err.Error(),
		}, err
	}
	if res.ModifiedCount == 0 {
		return &pb.UpdateTransactionResp{
			Status:  "error",
			Message: "Transaction not found",
		}, nil
	}

	return &pb.UpdateTransactionResp{
		Status:  "success",
		Message: "Transaction updated successfully",
	}, nil
}

func (repo *transactionRepositoryImpl) DeleteTransaction(ctx context.Context, request *pb.DeleteTransactionReq) (*pb.DeleteTransactionResp, error) {
	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "user_id", Value: request.UserId},
		{Key: "deleted_at", Value: time.Now()},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: nil}}}})

	if err != nil {
		return &pb.DeleteTransactionResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}

	if res.ModifiedCount == 0 {
		return &pb.DeleteTransactionResp{
			Status:  "error",
			Message: "transaction not found",
		}, nil
	}

	return &pb.DeleteTransactionResp{
		Status:  "success",
		Message: "Transaction deleted successfully",
	}, nil
}

func (repo *transactionRepositoryImpl) GetTransaction(ctx context.Context, request *pb.GetTransactionReq) (*pb.GetTransactionResp, error) {
	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "user_id", Value: request.UserId},
		{Key: "deleted_at", Value: nil},
	}

	var transaction models.GetTransaction
	err := repo.coll.FindOne(ctx, filter).Decode(&transaction)

	if err != nil {
		return nil, err
	}

	typeTrnsaction, err := enums.TypeTrnsactionParse(transaction.Type)
	if err != nil {
		return nil, err
	}

	return &pb.GetTransactionResp{
		Id:          transaction.Id,
		AccountId:   transaction.AccountId,
		UserId:      transaction.UserId,
		Type:        *typeTrnsaction,
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Date:        transaction.Date.Format("2006-01-02 15:04:05"),
	}, nil
}

func (repo *transactionRepositoryImpl) GetTransactionsList(ctx context.Context, request *pb.GetTransactionsListReq) (*pb.GetTransactionsListResp, error) {
	return nil, nil
}