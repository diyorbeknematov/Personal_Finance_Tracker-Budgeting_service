package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"budgeting-service/pkg/enums"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BudgetManagementRepo interface {
	CreateBudget(ctx context.Context, budget *pb.CreateBudgetReq) (*pb.CreateBudgetResp, error)
	UpdateBudget(ctx context.Context, budget *pb.UpdateBudgetReq) (*pb.UpdateBudgetResp, error)
	DeleteBudget(ctx context.Context, request *pb.DeleteBudgetReq) (*pb.DeleteBudgetResp, error)
	GetBudget(ctx context.Context, request *pb.GetBudgetReq) (*pb.GetBudgetResp, error)
	GetBudgetsList(ctx context.Context, request *pb.GetBudgetsReq) (*pb.GetBudgetsResp, error)
}

type budgetManagementRepoImpl struct {
	coll *mongo.Collection
}

func NewBudgetManagementRepo(db *mongo.Database) BudgetManagementRepo {
	return &budgetManagementRepoImpl{coll: db.Collection("budgets")}
}

func (repo *budgetManagementRepoImpl) CreateBudget(ctx context.Context, budget *pb.CreateBudgetReq) (*pb.CreateBudgetResp, error) {
	sdate, err := time.Parse("2006-01-02 15:04:05", budget.StartDate)
	if err != nil {
		return nil, err
	}
	edate, err := time.Parse("2006-01-02 15:04:05", budget.EndDate)
	if err != nil {
		return nil, err
	}
	_, err = repo.coll.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.NewString()},
		{Key: "user_id", Value: budget.UserId},
		{Key: "category_id", Value: budget.CategoryId},
		{Key: "amount", Value: budget.Amount},
		{Key: "period", Value: budget.Period.String()},
		{Key: "start_date", Value: sdate},
		{Key: "end_date", Value: edate},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
		{Key: "deleted_at", Value: nil},
	})

	if err != nil {
		return &pb.CreateBudgetResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}
	return &pb.CreateBudgetResp{
		Status:  "success",
		Message: "create budget successfully",
	}, nil
}

func (repo *budgetManagementRepoImpl) GetBudget(ctx context.Context, req *pb.GetBudgetReq) (*pb.GetBudgetResp, error) {
	filter := bson.D{
		{Key: "_id", Value: req.Id},
		{Key: "deleted_at", Value: nil},
	}

	var budget models.GetBudget
	err := repo.coll.FindOne(ctx, filter).Decode(&budget)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("budget not found")
	} else if err != nil {
		return nil, err
	}

	period, err := enums.PeriodParse(budget.Period)
	if err != nil {
		return nil, err
	}

	return &pb.GetBudgetResp{
		Id:         budget.ID,
		UserId:     budget.UserId,
		CategoryId: budget.CategoryId,
		Amount:     budget.Amount,
		Period:     *period,
		StartDate:  budget.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:    budget.EndDate.Format("2006-01-02 15:04:05"),
	}, nil
}

func (repo *budgetManagementRepoImpl) UpdateBudget(ctx context.Context, budget *pb.UpdateBudgetReq) (*pb.UpdateBudgetResp, error) {
	filter := bson.D{
		{Key: "_id", Value: budget.Id},
		{Key: "deleted_at", Value: nil},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "amount", Value: budget.Amount},
			{Key: "period", Value: budget.Period.String()},
			{Key: "start_date", Value: budget.StartDate},
			{Key: "end_date", Value: budget.EndDate},
			{Key: "updated_at", Value: time.Now()},
		}},
	})

	if err != nil {
		return &pb.UpdateBudgetResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}

	if res.MatchedCount == 0 {
		return &pb.UpdateBudgetResp{
			Status:  "error",
			Message: "budget not found",
		}, fmt.Errorf("budget not found")
	}

	return &pb.UpdateBudgetResp{
		Status:  "success",
		Message: "budget updated successfully",
	}, nil
}

func (repo *budgetManagementRepoImpl) DeleteBudget(ctx context.Context, req *pb.DeleteBudgetReq) (*pb.DeleteBudgetResp, error) {
	filter := bson.D{
		{Key: "_id", Value: req.Id},
		{Key: "deleted_at", Value: nil},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}},
	})

	if err != nil {
		return &pb.DeleteBudgetResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}

	if res.MatchedCount == 0 {
		return &pb.DeleteBudgetResp{
			Status:  "error",
			Message: "budget not found",
		}, fmt.Errorf("budget not found")
	}

	return &pb.DeleteBudgetResp{
		Status:  "success",
		Message: "deleted budget successfully",
	}, nil
}


func (repo *budgetManagementRepoImpl) GetBudgetsList(ctx context.Context, budget *pb.GetBudgetsReq) (*pb.GetBudgetsResp, error) {
	return nil, nil
}	