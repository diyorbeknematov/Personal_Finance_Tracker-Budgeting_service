package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGoal(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewGoalsRepository(db)
	goal := &pb.CreateGoalReq{
		UserId:       "test_user_id",
		Name:         "Test Goal",
		TargetAmount: 1000,
		Deadline:     "2022-12-31 00:00:00",
		Status:       "ACHIEVED",
	}

	resp, err := repo.CreateGoal(context.Background(), goal)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "success", resp.Status)
}

func TestUpdateGoal(t *testing.T) {
	db, err := ConnectToMongoDB()
    if err != nil {
        t.Fatal(err)
    }
    defer db.Client().Disconnect(context.Background())

    repo := NewGoalsRepository(db)
    goal := &pb.UpdateGoalReq{
        Id:          "f5764263-bddf-4289-b4b1-ecc2e5e0931a",
        Name:         "Updated Test Goal",
        TargetAmount: 2000,
        Deadline:     "2023-01-31",
        Status:       "INPORGRESS",
    }

    resp, err := repo.UpdateGoal(context.Background(), goal)
    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, "success", resp.Status)
}

func TestDeleteGoal(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err!= nil {
        t.Fatal(err)
    }
	defer db.Client().Disconnect(context.Background())
	repo := NewGoalsRepository(db)

	request := &pb.DeleteGoalReq{
        Id: "f5764263-bddf-4289-b4b1-ecc2e5e0931a",
    }
	resp, err := repo.DeleteGoal(context.Background(), request)
	if err!= nil {
        t.Fatal(err)
    }
	assert.Equal(t, "success", resp.Status)
}

func TestGetGoal(t *testing.T) {
	db, err := ConnectToMongoDB()
    if err!= nil {
        t.Fatal(err)
    }
    defer db.Client().Disconnect(context.Background())
    repo := NewGoalsRepository(db)

    request := &pb.GetGoalReq{
        Id: "5e8ad7a5-fd85-44dd-bb7a-3087a5b9c7e5",
    }
    resp, err := repo.GetGoal(context.Background(), request)
    if err!= nil {
        t.Fatal(err)
    }
    assert.Equal(t, resp.Name, "Test Goal")
}
