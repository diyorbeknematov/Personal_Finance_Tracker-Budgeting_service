package service

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/storage"
	"context"
	"log/slog"
)

type ReportingNotifcation interface {
	GetSependingReport(ctx context.Context, request *pb.GetSependingReq) (*pb.GetSependingResp, error)
	GetIncomeReport(ctx context.Context, request *pb.GetIncomeReportReq) (*pb.GetIncomeReportResp, error)
	GetBudgetPerformance(ctx context.Context, request *pb.GetBudgetPerformanceReq) (*pb.GetBudgetPerformanceResp, error)
	GetGoalsProgress(ctx context.Context, request *pb.GetGoalProgressReq) (*pb.GetGoalProgressResp, error)

	SendNotification(context.Context, *pb.SendNotificationReq) (*pb.SendNotificationResp, error)
	GetNotificationList(context.Context, *pb.GetNotificationsListReq) (*pb.GetNotificationsListResp, error)
	GetNotification(context.Context, *pb.GetNotificationReq) (*pb.GetNotificationResp, error)
	UpdateNotification(context.Context, *pb.UpdateNotificationReq) (*pb.UpdateNotificationResp, error)
	DeleteNotification(context.Context, *pb.DeleteNotificationReq) (*pb.DeleteNotificationResp, error)
}

type reportingNotificationImpl struct {
	pb.UnimplementedReportingNotificationServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewReportingNotificationService(storage storage.IStorage, logger *slog.Logger) *reportingNotificationImpl {
	return &reportingNotificationImpl{
		storage: storage,
		logger:  logger,
	}
}

func (s *reportingNotificationImpl) GetSepending(ctx context.Context, request *pb.GetSependingReq) (*pb.GetSependingResp, error) {
	resp, err := s.storage.ReportingRepository().GetSependingReport(ctx, request)
	if err != nil {
		s.logger.Error("Get sepending report error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *reportingNotificationImpl) GetIncome(ctx context.Context, request *pb.GetIncomeReportReq) (*pb.GetIncomeReportResp, error) {
	resp, err := s.storage.ReportingRepository().GetIncomeReport(ctx, request)
	if err != nil {
		s.logger.Error("Get income report error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *reportingNotificationImpl) GetBudgetPerformance(ctx context.Context, request *pb.GetBudgetPerformanceReq) (*pb.GetBudgetPerformanceResp, error) {
	resp, err := s.storage.ReportingRepository().GetBudgetPerformance(ctx, request)
	if err != nil {
		s.logger.Error("Get budget performance report error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *reportingNotificationImpl) GoalProgress(ctx context.Context, request *pb.GetGoalProgressReq) (*pb.GetGoalProgressResp, error) {
	resp, err := s.storage.ReportingRepository().GetGoalsProgress(ctx, request)
	if err != nil {
		s.logger.Error("Get goals progress report error", "error", err)
		return resp, err
	}
	return resp, nil
}

func (s *reportingNotificationImpl) SendNotification(ctx context.Context, request *pb.SendNotificationReq) (*pb.SendNotificationResp, error) {
	resp, err := s.storage.NotificationRepository().SendNotification(ctx, request)
	if err!= nil {
        s.logger.Error("Send notification error", "error", err)
        return resp, err
    }
	return resp, nil
}

func (s *reportingNotificationImpl) GetNotificationList(ctx context.Context, request *pb.GetNotificationsListReq) (*pb.GetNotificationsListResp, error) {
	resp, err := s.storage.NotificationRepository().GetNotificationsList(ctx, request)
	if err!= nil {
        s.logger.Error("Get notification list error", "error", err)
        return resp, err
    }
	return resp, nil
}

func (s *reportingNotificationImpl) GetNotification(ctx context.Context, request *pb.GetNotificationReq) (*pb.GetNotificationResp, error) {
	resp, err := s.storage.NotificationRepository().GetNotification(ctx, request)
	if err!= nil {
        s.logger.Error("Get notification error", "error", err)
        return resp, err
    }
	return resp, nil
}

func (s *reportingNotificationImpl) UpdateNotification(ctx context.Context, request *pb.UpdateNotificationReq) (*pb.UpdateNotificationResp, error) {
	resp, err := s.storage.NotificationRepository().UpdateNotification(ctx, request)
	if err != nil {
		s.logger.Error("Update notification error", "error", err)
        return resp, err
    }
	
	return resp, nil
}

func (s *reportingNotificationImpl) DeleteNotification(ctx context.Context, request *pb.DeleteNotificationReq) (*pb.DeleteNotificationResp, error) {
	resp, err := s.storage.NotificationRepository().DeleteNotification(ctx, request)
	if err!= nil {
        s.logger.Error("Delete notification error", "error", err)
        return resp, err
    }
	return resp, nil
}
