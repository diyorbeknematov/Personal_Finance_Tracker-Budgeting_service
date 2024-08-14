package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewCategoryRepository(db)
	category := &pb.CreateCategoryReq{
		UserId: "test_user_id",
		Name:   "Test Category",
		Type:   "EXPENSE",
	}

	resp, err := repo.CreateCategory(context.Background(), category)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
}

func TestUpdateCategory(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewCategoryRepository(db)
	category := &pb.UpdateCategoryReq{
		Id:   "d9d36573-5ee9-48e3-80a0-257172f64c72",
		Name: "Updated Test Category",
		Type: "INCOME",
	}

	resp, err := repo.UpdateCategory(context.Background(), category)
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
}

func TestDeleteCategory(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewCategoryRepository(db)
	resp, err := repo.DeleteCategory(context.Background(), &pb.DeleteCategoryReq{
		Id:     "d9d36573-5ee9-48e3-80a0-257172f64c72",
		UserId: "test_user_id",
	})
	assert.NoError(t, err)

	assert.Equal(t, "success", resp.Status)
}

func TestGetCategory(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewCategoryRepository(db)
	resp, err := repo.GetCategory(context.Background(), &pb.GetCategoryReq{
		Id: "76db8dc2-8e59-48d0-8f73-e09059008435",
	})
	assert.NoError(t, err)

	assert.Equal(t, "test_user_id", resp.UserId)
}

func TestGetCategoriesList(t *testing.T) {
	db, err := ConnectToMongoDB()
    if err!= nil {
        t.Fatal(err)
    }
    defer db.Client().Disconnect(context.Background())

    repo := NewCategoryRepository(db)
    resp, err := repo.GetCategoriesList(context.Background(), &pb.GetCategoriesReq{
        Name: "Test Category",
		Type: "EXPENSE",
    })
    assert.NoError(t, err)

    assert.NotNil(t, resp)
}