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
		{Key: "period", Value: budget.Period},
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

	return &pb.GetBudgetResp{
		Id:         budget.ID,
		UserId:     budget.UserId,
		CategoryId: budget.CategoryId,
		Amount:     budget.Amount,
		Period:     budget.Period,
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
			{Key: "period", Value: budget.Period},
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
	pipeline := createBudgetFilters(budget)
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

	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (budget.Page - 1) * budget.Limit}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: budget.Limit}})

	cursor, err := repo.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var budgets []*pb.Budget
	for cursor.Next(ctx) {
		var budget models.GetBudget
		err := cursor.Decode(&budget)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, &pb.Budget{
			Id:         budget.ID,
			UserId:     budget.UserId,
			CategoryId: budget.CategoryId,
			Amount:     budget.Amount,
			Period:     budget.Period,
			StartDate:  budget.StartDate.Format("2006-01-02 15:04:05"),
			EndDate:    budget.EndDate.Format("2006-01-02 15:04:05"),
		})
	}
	if err := cursor.Close(ctx); err != nil {
		return nil, err
	}
	if len(budgets) == 0 {
		return nil, fmt.Errorf("budgets not found")
	}

	return &pb.GetBudgetsResp{
		Budgets:    budgets,
		TotalCount: int64(totalCount),
	}, nil
}

func createBudgetFilters(request *pb.GetBudgetsReq) mongo.Pipeline {
	var pipeline mongo.Pipeline

	// Foydalanuvchi ID bo'yicha filter
	if request.UserId != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: request.UserId},
			}},
		})
	}

	// Kategoriya nomi bo'yicha filter
	if request.CategoreName != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "categories"},
				{Key: "localField", Value: "category_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "category"},
			}},
		})
		pipeline = append(pipeline, bson.D{
			{Key: "$unwind", Value: "$category"},
		})
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "category.name", Value: bson.D{
					{Key: "$regex", Value: ".*" + request.CategoreName + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}

	// Miqdor bo'yicha filter
	if request.Amount != 0 {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "amount", Value: bson.D{
					{Key: "$gte", Value: request.Amount},
				}},
			}},
		})
	}

	// Davr bo'yicha filter (period)
	if request.Period != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "period", Value: bson.D{
					{Key: "$regex", Value: ".*" + request.Period + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}

	if request.StartDate != "" && request.EndDate != "" {
		startDate, err := time.Parse(time.RFC3339, request.StartDate)
		if err != nil {
			return pipeline
		}
		endDate, err := time.Parse(time.RFC3339, request.EndDate)
		if err != nil {
			// Handle error
			return pipeline
		}

		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "start_date", Value: bson.D{
					{Key: "$gte", Value: startDate},
				}},
				{Key: "end_date", Value: bson.D{
					{Key: "$lte", Value: endDate},
				}},
			}},
		})
	}

	return pipeline
}
