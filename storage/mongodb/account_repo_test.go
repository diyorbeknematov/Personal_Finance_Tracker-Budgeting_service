package mongodb

import (
	"budgeting-service/config"
	pb "budgeting-service/generated/budgeting"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToMongo(t *testing.T) {
	cfg := config.Load()
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())
	t.Log("Connected to MongoDB successfully")

	assert.Equal(t, cfg.MONGODB_NAME, db.Name())
}

func TestCreateAccount(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewAccountRepository(db)
	account := &pb.CreateAccountReq{
		UserId:   "test_user_id",
		Name:     "Test Account",
		Type:     "CHECKING",
		Balance:  1000.0,
		Currency: "USD",
	}
	resp, err := repo.CreateAccount(context.Background(), account)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "success", resp.Status)
}

func TestUpdateAccount(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewAccountRepository(db)
	account := &pb.UpdateAccountReq{
		Id:       "1e891746-d2bb-44d0-bab0-4660ee52d12d",
		Name:     "Updated Test Account",
		Type:     "SAVINGS",
		Balance:  2000.0,
		Currency: "USD",
	}
	resp, err := repo.UpdateAccount(context.Background(), account)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "success", resp.Status)
}

func TestDeleteAccount(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewAccountRepository(db)
	resp, err := repo.DeleteAccount(context.Background(), &pb.DeleteAccountReq{
		Id:     "test_account_id",
		UserId: "test_user_id",
	})
	if err != nil {
		t.Fatal(err)
	}

	if assert.Equal(t, "account not found", resp.Message) {
		return
	}
	assert.Equal(t, "success", resp.Status)
}

func TestGetAccountsList(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewAccountRepository(db)
	resp, err := repo.GetAccountsList(context.Background(), &pb.GetAccountsListReq{
		UserId: "test_user_id",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(resp.Accounts))
}

func TestGetAccount(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())

	repo := NewAccountRepository(db)
	resp, err := repo.GetAccount(context.Background(), &pb.GetAccountReq{
		Id:     "116ed66a-d57f-4fd1-b2ef-1310b5367008",
	})

	if assert.Nil(t, err) {
		t.Log("GetAccount returned an error")
	} else {
		assert.Equal(t, "ddb813cb-28bd-4edc-819c-e492f7eeaa5c", resp.GetUserId())
	}
}
