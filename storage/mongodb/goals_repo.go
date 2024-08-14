package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoalsRepository interface {
	CreateGoal(ctx context.Context, goal *pb.CreateGoalReq) (*pb.CreateGoalResp, error)
	UpdateGoal(ctx context.Context, goal *pb.UpdateGoalReq) (*pb.UpdateGoalResp, error)
	DeleteGoal(ctx context.Context, request *pb.DeleteGoalReq) (*pb.DeleteGoalResp, error)
	GetGoal(ctx context.Context, request *pb.GetGoalReq) (*pb.GetGoalResp, error)
	GetGoalsList(ctx context.Context, request *pb.GetGoalsReq) (*pb.GetGoalsResp, error)
}

type goalsRepositoryImpl struct {
	coll *mongo.Collection
}

func NewGoalsRepository(db *mongo.Database) GoalsRepository {
	return &goalsRepositoryImpl{coll: db.Collection("goals")}
}

func (repo *goalsRepositoryImpl) CreateGoal(ctx context.Context, goal *pb.CreateGoalReq) (*pb.CreateGoalResp, error) {
	deadline, err := time.Parse("2006-01-02 15:04:05", goal.Deadline)
	if err != nil {
		return nil, err
	}
	_, err = repo.coll.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.NewString()},
		{Key: "user_id", Value: goal.UserId},
		{Key: "name", Value: goal.Name},
		{Key: "target_amount", Value: goal.TargetAmount},
		{Key: "current_amount", Value: 0},
		{Key: "deadline", Value: deadline},
		{Key: "status", Value: goal.Status},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
		{Key: "deleted_at", Value: nil},
	})

	if err != nil {
		return &pb.CreateGoalResp{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}

	return &pb.CreateGoalResp{
		Status:  "success",
		Message: "Goal created successfully",
	}, nil
}

func (repo *goalsRepositoryImpl) UpdateGoal(ctx context.Context, goal *pb.UpdateGoalReq) (*pb.UpdateGoalResp, error) {
	filter := bson.D{
		{Key: "_id", Value: goal.Id},
		{Key: "deleted_at", Value: nil},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: goal.Name},
			{Key: "target_amount", Value: goal.TargetAmount},
			{Key: "deadline", Value: goal.Deadline},
			{Key: "status", Value: goal.Status},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return &pb.UpdateGoalResp{
			Status:  "error",
			Message: "Error updating goal: " + err.Error(),
		}, nil
	}

	if res.ModifiedCount == 0 {
		return &pb.UpdateGoalResp{
			Status:  "error",
			Message: "Goal not found",
		}, nil
	}

	return &pb.UpdateGoalResp{
		Status:  "success",
		Message: "Goal updated successfully",
	}, nil
}

func (repo *goalsRepositoryImpl) DeleteGoal(ctx context.Context, request *pb.DeleteGoalReq) (*pb.DeleteGoalResp, error) {
	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "deleted_at", Value: nil},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}},
	})

	if err != nil {
		return &pb.DeleteGoalResp{
			Status:  "error",
			Message: "Error deleting goal: " + err.Error(),
		}, nil
	}

	if res.ModifiedCount == 0 {
		return &pb.DeleteGoalResp{
			Status:  "error",
			Message: "Goal not found",
		}, nil
	}

	return &pb.DeleteGoalResp{
		Status:  "success",
		Message: "Goal deleted successfully",
	}, nil
}

func (repo *goalsRepositoryImpl) GetGoal(ctx context.Context, request *pb.GetGoalReq) (*pb.GetGoalResp, error) {
	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "deleted_at", Value: nil},
	}

	goal := models.GetGoal{}
	err := repo.coll.FindOne(ctx, filter).Decode(&goal)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("goal not found")
	} else if err != nil {
		return nil, err
	}

	fmt.Println(goal)
	return &pb.GetGoalResp{
		Id:            goal.ID,
		UserId:        goal.UserId,
		Name:          goal.Name,
		TargetAmount:  goal.TargetAmount,
		CurrentAmount: goal.CurrentAmount,
		Deadline:      goal.Deadline.Format("2006-01-02 15:04:05"),
		Status:        goal.Status,
	}, nil
}

func (repo *goalsRepositoryImpl) GetGoalsList(ctx context.Context, request *pb.GetGoalsReq) (*pb.GetGoalsResp, error) {
	pipeline := createGoalFilters(request)

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

	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (request.Page - 1) * request.Limit}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: request.Limit}})

	cursor, err := repo.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var goals []*pb.Goal
	for cursor.Next(ctx) {
		var goal models.GetGoal
		err := cursor.Decode(&goal)
		if err != nil {
			return nil, err
		}

		goals = append(goals, &pb.Goal{
			Id:            goal.ID,
			UserId:        goal.UserId,
			Name:          goal.Name,
			TargetAmount:  goal.TargetAmount,
			CurrentAmount: goal.CurrentAmount,
			Deadline:      goal.Deadline.Format("2006-01-02 15:04:05"),
			Status:        goal.Status,
		})
	}

	if len(goals) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &pb.GetGoalsResp{
		Goals:      goals,
		TotalCount: int64(totalCount),
		Page:       request.Page,
		Limit:      request.Limit,
	}, nil
}

func createGoalFilters(request *pb.GetGoalsReq) mongo.Pipeline {
	pipeline := mongo.Pipeline{}
	if request.UserId != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
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
	if request.TargetAmount != 0 {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "target_amount", Value: request.TargetAmount},
			}},
		})
	}
	if request.Deadline != "" {
		deadline, err := time.Parse("2006-01-02 15:04:05", request.Deadline)
		if err == nil {
			pipeline = append(pipeline, bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "deadline", Value: bson.D{
						{Key: "$gte", Value: deadline},
						{Key: "$lte", Value: time.Now()},
					}},
				}},
			})
		}
	}

	if request.Status != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "status", Value: request.Status},
			}},
		})
	}

	return pipeline
}
