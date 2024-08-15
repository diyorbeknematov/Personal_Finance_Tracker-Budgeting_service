package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
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
		{Key: "type", Value: transaction.Type},
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
		{Key: "deleted_at", Value: nil},
	}

	var transaction models.GetTransaction
	err := repo.coll.FindOne(ctx, filter).Decode(&transaction)

	if err != nil {
		return nil, err
	}

	return &pb.GetTransactionResp{
		Id:          transaction.Id,
		AccountId:   transaction.AccountId,
		UserId:      transaction.UserId,
		Type:        transaction.Type,
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Date:        transaction.Date.Format("2006-01-02 15:04:05"),
	}, nil
}

func (repo *transactionRepositoryImpl) GetTransactionsList(ctx context.Context, request *pb.GetTransactionsListReq) (*pb.GetTransactionsListResp, error) {
	pipeline := createTransactionFilter(request)
	// Hujjatlarni sanash uchun `pipeline` ni `Aggregate` bilan ishlating
	countPipeline := append(pipeline, bson.D{{Key: "$count", Value: "totalCount"}})
	countCursor, err := repo.coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, err
	}
	var countResult []bson.M
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, err
	}
	totalCount := int32(0)
	if len(countResult) > 0 {
		totalCount = countResult[0]["totalCount"].(int32)
	}

	// Pagination uchun `skip` va `limit` qo'shing
	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (request.Page - 1) * request.Limit}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: request.Limit}})

	cursor, err := repo.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*pb.Transaction
	for cursor.Next(ctx) {
		var transaction models.GetTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &pb.Transaction{
			Id:          transaction.Id,
			AccountId:   transaction.AccountId,
			UserId:      transaction.UserId,
			Type:        transaction.Type,
			Amount:      transaction.Amount,
			Description: transaction.Description,
			Date:        transaction.Date.Format("2006-01-02 15:04:05"),
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &pb.GetTransactionsListResp{
		TotalCount:   int64(totalCount),
		Transactions: transactions,
		Page:         request.Page,
		Limit:        request.Limit,
	}, nil
}

func createTransactionFilter(request *pb.GetTransactionsListReq) mongo.Pipeline {
	var pipeline mongo.Pipeline
	if request.UserId != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: request.UserId}}}})
	}
	if request.AccountName != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "accounts"},
				{Key: "localField", Value: "account_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "account"},
			}},
		})
		pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: "$account"}})
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "account.name", Value: bson.D{
					{Key: "$regex", Value: ".*" + request.AccountName + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}
	if request.CategoryName != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "categories"},
				{Key: "localField", Value: "category_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "category"},
			}},
		})
		pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: "$category"}})
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "category.name", Value: bson.D{
					{Key: "$regex", Value: ".*" + request.CategoryName + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}
	if request.Amount != 0 {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "amount", Value: bson.D{
					{Key: "$gte", Value: request.Amount},
				}},
			}},
		})
	}
	if request.Description != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "description", Value: bson.D{
					{Key: "$regex", Value: ".*" + request.Description + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}

	if request.Type != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "type", Value: bson.D{
					{Key: "$regex", Value: ".*" + request.Type + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}

	pipeline = append(pipeline, bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "delete_at", Value: nil},
		}},
	})

	return pipeline
}
