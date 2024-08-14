package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateBudget(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewBudgetManagementRepo(db)
	budget := &pb.CreateBudgetReq{
		UserId:     "test_user_id",
		CategoryId: "test_category_id",
		Amount:     1000,
		Period:     "MONTHLY",
		StartDate:  time.Now().Format("2006-01-02 15:04:05"),
		EndDate:    time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02 15:04:05"),
	}

	resp, err := repo.CreateBudget(context.Background(), budget)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
}

func TestGetBudget(t *testing.T) {
	db, err := ConnectToMongoDB()
    if err!= nil {
        t.Fatal(err)
    }
    defer db.Client().Disconnect(context.Background())

    repo := NewBudgetManagementRepo(db)
    req := &pb.GetBudgetReq{
        Id: "b2828dd7-9460-448e-bd43-c754ad775c01",
    }

    resp, err := repo.GetBudget(context.Background(), req)
    assert.NoError(t, err)
    assert.NotNil(t, resp)
}

func TestUpdateBudget(t *testing.T) {
	db, err := ConnectToMongoDB()
    if err!= nil {
        t.Fatal(err)
    }
    defer db.Client().Disconnect(context.Background())

    repo := NewBudgetManagementRepo(db)
    budget := &pb.UpdateBudgetReq{
        Id:   "b2828dd7-9460-448e-bd43-c754ad775c01",
        CategoryId: "test_category_id",
        Amount:     1000,
        Period:     "MONTHLY",
        StartDate:  time.Now().Format("2006-01-02 15:04:05"),
        EndDate:    time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02 15:04:05"),
	}
	resp, err := repo.UpdateBudget(context.Background(), budget)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
}

func TestDeleteBudget(t *testing.T) {
	db, err := ConnectToMongoDB()
    if err!= nil {
        t.Fatal(err)
    }
    defer db.Client().Disconnect(context.Background())

    repo := NewBudgetManagementRepo(db)
    req := &pb.DeleteBudgetReq{
        Id: "b2828dd7-9460-448e-bd43-c754ad775c01",
		UserId: "test_user_id",
    }

    resp, err := repo.DeleteBudget(context.Background(), req)
    assert.NoError(t, err)
    assert.Equal(t, "success", resp.Status)
}