package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// <--- Transactions and budgets are not tested in this test suite, as they are specific to the budget
func TestCreateTransaction(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewTransactionRepository(db)
	transaction := &pb.CreateTransactionReq{
		UserId:      "test_user_id",
		AccountId:   "test_account_id",
		Amount:      100.0,
		Type:        "expence",
		Description: "Test Transaction",
		CategoryId:  "test_category",
		Date:        "2022-01-01 00:00:00",
	}
	resp, err := repo.CreateTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "success", resp.Status)
}

func TestUpdateTransaction(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewTransactionRepository(db)
	transaction := &pb.UpdateTransactionReq{
		Id:          "1e891746-d2bb-44d0-bab0-4660ee52d12d",
		Amount:      200.0,
		Type:        "income",
		Description: "Updated Test Transaction",
		Date:        "2022-01-02 00:00:00",
	}

	resp, err := repo.UpdateTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "success", resp.Status)
}

func TestDeleteTransaction(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewTransactionRepository(db)
	resp, err := repo.DeleteTransaction(context.Background(), &pb.DeleteTransactionReq{
		Id:     "test_transaction_id",
		UserId: "test_user_id",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "transaction not found", resp.Message)
	assert.Equal(t, "success", resp.Status)
}

func TestGetTransactions(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewTransactionRepository(db)
	resp, err := repo.GetTransactionsList(context.Background(), &pb.GetTransactionsListReq{
		Type: "income",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.TotalCount)
}

func TestGetTransaction(t *testing.T) {
	db, err := ConnectToMongoDB()
    if err!= nil {
        t.Fatal(err)
    }
    defer db.Client().Disconnect(context.Background())

    repo := NewTransactionRepository(db)
    resp, err := repo.GetTransaction(context.Background(), &pb.GetTransactionReq{
        Id:     "70ea539d-3099-40fb-bbd9-155f05d9d5e8",
    })
    if err!= nil {
        t.Fatal(err)
    }
    assert.NotNil(t, resp)
    assert.Equal(t, "test_user_id", resp.UserId)
}