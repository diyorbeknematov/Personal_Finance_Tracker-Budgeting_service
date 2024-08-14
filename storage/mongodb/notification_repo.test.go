package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendNotification(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())
	repo := NewNotificationRepository(db)
	resp, err := repo.SendNotification(context.Background(), &pb.SendNotificationReq{
		UserId:  "test_user_id",
		Message: "Test notification",
		Type:    "test_type",
		Status:  "test_status",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.Status, "success")
}

func TestGetNotification(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())
	repo := NewNotificationRepository(db)
	resp, err := repo.GetNotification(context.Background(), &pb.GetNotificationReq{
		Id: "test_notification_id",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.Status, "success")
}

func TestDeleteNotification(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())
	repo := NewNotificationRepository(db)
	resp, err := repo.DeleteNotification(context.Background(), &pb.DeleteNotificationReq{
		Id: "test_notification_id",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resp.Status, "success")
}

func TestGetNotificationsList(t *testing.T) {
	db, err := ConnectToMongoDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Client().Disconnect(context.Background())
	repo := NewNotificationRepository(db)
	resp, err := repo.GetNotificationsList(context.Background(), &pb.GetNotificationsListReq{
		UserId: "test_user_id",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, resp)
}
