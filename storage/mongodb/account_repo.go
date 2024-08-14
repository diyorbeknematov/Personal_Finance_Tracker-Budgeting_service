package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *pb.CreateAccountReq) (*pb.CreateAccountResp, error)
	UpdateAccount(ctx context.Context, account *pb.UpdateAccountReq) (*pb.UpdateAccountResp, error)
	DeleteAccount(ctx context.Context, request *pb.DeleteAccountReq) (*pb.DeleteAccountResp, error)
	GetAccount(ctx context.Context, request *pb.GetAccountReq) (*pb.GetAccountResp, error)
	GetAccountsList(ctx context.Context, request *pb.GetAccountsListReq) (*pb.GetAccountsListResp, error)
}

type accountRepositoryImpl struct {
	coll *mongo.Collection
}

func NewAccountRepository(db *mongo.Database) AccountRepository {
	return &accountRepositoryImpl{coll: db.Collection("accounts")}
}

func (repo *accountRepositoryImpl) CreateAccount(ctx context.Context, account *pb.CreateAccountReq) (*pb.CreateAccountResp, error) {
	_, err := repo.coll.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.NewString()},
		{Key: "user_id", Value: account.UserId},
		{Key: "name", Value: account.Name},
		{Key: "type", Value: account.Type},
		{Key: "balance", Value: account.Balance},
		{Key: "currency", Value: account.Currency},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
		{Key: "deleted_at", Value: nil},
	})

	if err != nil {
		return &pb.CreateAccountResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}
	return &pb.CreateAccountResp{
		Status:  "success",
		Message: "created account successfully",
	}, nil
}

func (repo *accountRepositoryImpl) UpdateAccount(ctx context.Context, account *pb.UpdateAccountReq) (*pb.UpdateAccountResp, error) {
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: account.Name},
			{Key: "type", Value: account.Type},
			{Key: "balance", Value: account.Balance},
			{Key: "currency", Value: account.Currency},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	filter := bson.D{
		{Key: "_id", Value: account.Id},
		{Key: "deleted_at", Value: nil},
	}

	_, err := repo.coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return &pb.UpdateAccountResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}
	return &pb.UpdateAccountResp{
		Status:  "success",
		Message: "updated account successfully",
	}, nil
}

func (repo *accountRepositoryImpl) DeleteAccount(ctx context.Context, request *pb.DeleteAccountReq) (*pb.DeleteAccountResp, error) {
	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "user_id", Value: request.UserId},
		{Key: "deleted_at", Value: time.Now()},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: nil}}}})

	if err != nil {
		return &pb.DeleteAccountResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}

	if res.ModifiedCount == 0 {
		return &pb.DeleteAccountResp{
			Status:  "error",
			Message: "account not found",
		}, nil
	}

	return &pb.DeleteAccountResp{
		Status:  "success",
		Message: "deleted account successfully",
	}, nil
}

func (repo *accountRepositoryImpl) GetAccount(ctx context.Context, request *pb.GetAccountReq) (*pb.GetAccountResp, error) {
	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "user_id", Value: request.UserId},
		{Key: "deleted_at", Value: nil},
	}

	var account models.GetAccount
	err := repo.coll.FindOne(ctx, filter).Decode(&account)

	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResp{
		Id:       account.ID,
		UserId:   account.UserId,
		Name:     account.Name,
		Type:     account.Type,
		Balance:  account.Balance,
		Currency: account.Currency,
	}, nil
}

func (repo *accountRepositoryImpl) GetAccountsList(ctx context.Context, request *pb.GetAccountsListReq) (*pb.GetAccountsListResp, error) {
	pipeline := createFilters(request)

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

	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (request.Offset - 1) * request.Limit}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: request.Limit}})

	cursor, err := repo.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var accounts []*pb.Account
	for cursor.Next(ctx) {
		var account pb.Account
		err := cursor.Decode(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	if err := cursor.Close(ctx); err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("account not found")
	}
	return &pb.GetAccountsListResp{
		Limit:      request.Limit,
		Offset:     request.Offset,
		TotalCount: int64(totalCount),
		Accounts:   accounts,
	}, nil
}

func createFilters(request *pb.GetAccountsListReq) mongo.Pipeline {
	pipeline := mongo.Pipeline{}
	if request.UserId != "" {
		pipeline = append(pipeline, bson.D{{
			Key: "$match", Value: bson.D{
				{Key: "user_id", Value: request.UserId},
			}},
		})
	}
	if request.Name != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "name", Value: bson.D{
					{Key: "$regex", Value: ".*" + request.Name + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}

	if request.Type != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "type", Value: request.Type},
			}},
		})
	}

	if request.Currency != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "currency", Value: bson.D{
					{Key: "$regex", Value: ".*" + request.Currency + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}

	if request.Balance != 0 {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "balance", Value: bson.D{
					{Key: "$gte", Value: request.Balance},
				}},
			}},
		})
	}

	if request.CreatedAt {
		now := time.Now()
		oneMonthAgo := now.AddDate(0, -1, 0)
		pipeline = append(pipeline, bson.D{{
			Key: "$match", Value: bson.D{
				{Key: "created_at", Value: bson.D{
					{Key: "$gte", Value: oneMonthAgo},
				}},
			}},
		})
	}

	pipeline = append(pipeline, bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "deleted_at", Value: nil},
		}},
	})

	return pipeline
}

// < --- END OF Account  Implementation --- >
