package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository interface {
	SendNotification(ctx context.Context, notification *pb.SendNotificationReq) (*pb.SendNotificationResp, error)
	GetNotification(ctx context.Context, notification *pb.GetNotificationReq) (*pb.GetNotificationResp, error)
	GetNotificationsList(ctx context.Context, request *pb.GetNotificationsListReq) (*pb.GetNotificationsListResp, error)
	DeleteNotification(ctx context.Context, notification *pb.DeleteNotificationReq) (*pb.DeleteNotificationResp, error)
	UpdateNotification(ctx context.Context, notification *pb.UpdateNotificationReq) (*pb.UpdateNotificationResp, error)
}

type notificationRepositoryImpl struct {
	coll *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &notificationRepositoryImpl{coll: db.Collection("notifications")}
}

func (repo *notificationRepositoryImpl) SendNotification(ctx context.Context, notification *pb.SendNotificationReq) (*pb.SendNotificationResp, error) {
	_, err := repo.coll.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.NewString()},
		{Key: "user_id", Value: notification.UserId},
		{Key: "message", Value: notification.Message},
		{Key: "type", Value: notification.Type},
		{Key: "status", Value: notification.Status},
		{Key: "is_read", Value: false},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
		{Key: "deleted_at", Value: nil},
	})
	if err != nil {
		return &pb.SendNotificationResp{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}
	return &pb.SendNotificationResp{
		Status:  "success",
		Message: "Notification sent successfully",
	}, nil
}

func (repo *notificationRepositoryImpl) GetNotification(ctx context.Context, notification *pb.GetNotificationReq) (*pb.GetNotificationResp, error) {
	filter := bson.D{
		{Key: "_id", Value: notification.Id},
		{Key: "deleted_at", Value: nil},
	}

	var notifcation models.Notification
	err := repo.coll.FindOne(ctx, filter).Decode(&notifcation)

	if err != nil {
		return nil, err
	}
	return &pb.GetNotificationResp{
		Id:        notifcation.ID,
		Message:   notifcation.Message,
		Type:      notifcation.Type,
		Status:    notifcation.Status,
		CreatedAt: notifcation.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (repo *notificationRepositoryImpl) DeleteNotification(ctx context.Context, notification *pb.DeleteNotificationReq) (*pb.DeleteNotificationResp, error) {
	filter := bson.D{
		{Key: "_id", Value: notification.Id},
		{Key: "deleted_at", Value: nil},
	}
	_, err := repo.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "deleted_at", Value: time.Now()}}}})
	if err != nil {
		return &pb.DeleteNotificationResp{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}
	return &pb.DeleteNotificationResp{
		Status:  "success",
		Message: "Notification deleted successfully",
	}, nil
}

func (repo *notificationRepositoryImpl) UpdateNotification(ctx context.Context, notification *pb.UpdateNotificationReq) (*pb.UpdateNotificationResp, error) {
	filter := bson.D{
		{Key: "_id", Value: notification.Id},
		{Key: "deleted_at", Value: nil},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "is_read", Value: notification.IsRead},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	_, err := repo.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return &pb.UpdateNotificationResp{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}
	return &pb.UpdateNotificationResp{
		Status:  "success",
		Message: "Notification updated successfully",
	}, nil
}

func (repo *notificationRepositoryImpl) GetNotificationsList(ctx context.Context, notification *pb.GetNotificationsListReq) (*pb.GetNotificationsListResp, error) {
	filter := bson.D{
		{Key: "user_id", Value: notification.UserId},
		{Key: "deleted_at", Value: nil},
	}

	var notifcations []*pb.Notification
	cursor, err := repo.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var notifcation models.Notification
		err := cursor.Decode(&notifcation)
		if err != nil {
			return nil, err
		}
		notifcations = append(notifcations, &pb.Notification{
			Id:        notifcation.ID,
			Message:   notifcation.Message,
			Type:      notifcation.Type,
			Status:    notifcation.Status,
			IsRead:    notifcation.IsRead,
			CreatedAt: notifcation.CreatedAt.Format("2006-01-02 15:04:05"),
		})

	}

	return &pb.GetNotificationsListResp{
		NotificationList: notifcations,
	}, nil
}
