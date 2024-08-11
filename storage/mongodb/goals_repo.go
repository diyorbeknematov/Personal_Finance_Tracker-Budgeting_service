package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"budgeting-service/pkg/enums"
	"context"
	"errors"
	"fmt"
	"log"
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
		{Key: "userId", Value: goal.UserId},
		{Key: "name", Value: goal.Name},
		{Key: "target_amount", Value: goal.TargetAmount},
		{Key: "current_amount", Value: 0},
		{Key: "deadline", Value: deadline},
		{Key: "status", Value: goal.Status.String()},
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
			{Key: "status", Value: goal.Status.String()},
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

	goalStatus, err := enums.GoalStatusParse(goal.Status)
	if err != nil {
		log.Println("Error parsing goal status: " + err.Error())
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
		Status:        *goalStatus,
	}, nil
}

func (repo *goalsRepositoryImpl) GetGoalsList(ctx context.Context, request *pb.GetGoalsReq) (*pb.GetGoalsResp, error) {
	filter := bson.D{
		{Key: "userId", Value: request.UserId},
		{Key: "deleted_at", Value: nil},
	}

	cursor, err := repo.coll.Find(ctx, filter)
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

		goalStatus, err := enums.GoalStatusParse(goal.Status)
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
			Status:        *goalStatus,
		})
	}

	return &pb.GetGoalsResp{
		Goals: goals,
	}, nil
}
