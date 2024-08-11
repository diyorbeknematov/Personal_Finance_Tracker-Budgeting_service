package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"context"
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

	var account pb.GetAccountResp
	err := repo.coll.FindOne(ctx, filter).Decode(&account)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (repo *accountRepositoryImpl) GetAccountsList(ctx context.Context, request *pb.GetAccountsListReq) (*pb.GetAccountsListResp, error) {
	return nil, nil
}

// < --- END OF Account  Implementation --- >
