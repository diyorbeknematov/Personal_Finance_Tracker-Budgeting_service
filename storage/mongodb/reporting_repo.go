package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportingRepository interface {
}

type reportingRepositoryImpl struct {
	db *mongo.Database
}

func NewReportingRepository(db *mongo.Database) ReportingRepository {
	return &reportingRepositoryImpl{db: db}
}

func (repo *reportingRepositoryImpl) GetSependingReport(ctx context.Context, request *pb.GetSependingReq) (*pb.GetSependingResp, error) {
	pipeline := mongo.Pipeline{}

	// Boshlang'ich match
	pipeline = append(pipeline, bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: request.UserId},
			{Key: "type", Value: pb.TypeTrnsaction_EXPENSE.String()},
		}},
	})

	// Vaqt oralig'i uchun match
	if request.GetYearly() {
		// Oxirgi 1 yil uchun vaqt oralig'i
		startOfYear := time.Now().AddDate(-1, 0, 0)
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "date", Value: bson.D{
					{Key: "$gte", Value: startOfYear},
				}},
			}},
		})
	} else if request.GetMonthly() {
		// Oxirgi 1 oy uchun vaqt oralig'i
		startOfMonth := time.Now().AddDate(0, -1, 0)
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "date", Value: bson.D{
					{Key: "$gte", Value: startOfMonth},
				}},
			}},
		})
	}

	// Vaqt bo'yicha agregatsiya
	if request.GetYearly() {
		pipeline = append(pipeline, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_amount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			}},
		})
	} else if request.GetMonthly() {
		pipeline = append(pipeline, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$dateToString", Value: bson.D{
						{Key: "format", Value: "%Y-%m"}, // Yil va oy
						{Key: "date", Value: "$date"},
					}},
				}},
				{Key: "total_amount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			}},
		})
	} else {
		// Agar vaqt intervali ko'rsatilmagan bo'lsa
		pipeline = append(pipeline, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_amount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			}},
		})
	}

	// Aggregatsiya
	cursor, err := repo.db.Collection("transactions").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // Ensure cursor is closed after usage

	var result pb.GetSependingResp
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	} else {
		if err := cursor.Err(); err != nil {
			return nil, err
		}
		return nil, mongo.ErrNoDocuments
	}

	// Moslashtirilgan javob
	result.Yearly = request.GetYearly()
	result.Monthly = request.GetMonthly()

	return &result, nil
}

func (repo *reportingRepositoryImpl) GetIncomeReport(ctx context.Context, request *pb.GetIncomeReportReq) (*pb.GetIncomeReportResp, error) {
	pipeline := mongo.Pipeline{}

	// Boshlang'ich match
	pipeline = append(pipeline, bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: request.UserId},
			{Key: "type", Value: pb.TypeTrnsaction_INCOME.String()},
		}},
	})

	// Vaqt oralig'i uchun match
	if request.GetYearly() {
		// Oxirgi 1 yil uchun vaqt oralig'i
		startOfYear := time.Now().AddDate(-1, 0, 0)
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "date", Value: bson.D{
					{Key: "$gte", Value: startOfYear},
				}},
			}},
		})
	} else if request.GetMonthly() {
		// Oxirgi 1 oy uchun vaqt oralig'i
		startOfMonth := time.Now().AddDate(0, -1, 0)
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "date", Value: bson.D{
					{Key: "$gte", Value: startOfMonth},
				}},
			}},
		})
	}

	// Vaqt bo'yicha agregatsiya
	if request.GetYearly() {
		pipeline = append(pipeline, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_amount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			}},
		})
	} else if request.GetMonthly() {
		pipeline = append(pipeline, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$dateToString", Value: bson.D{
						{Key: "format", Value: "%Y-%m"}, // Yil va oy
						{Key: "date", Value: "$date"},
					}},
				}},
				{Key: "total_amount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			}},
		})
	} else {
		// Agar vaqt intervali ko'rsatilmagan bo'lsa
		pipeline = append(pipeline, bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_amount", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			}},
		})
	}

	// Aggregatsiya
	cursor, err := repo.db.Collection("transactions").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // Ensure cursor is closed after usage

	var result pb.GetIncomeReportResp
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	} else {
		if err := cursor.Err(); err != nil {
			return nil, err
		}
		return nil, mongo.ErrNoDocuments
	}

	// Moslashtirilgan javob
	result.Yearly = request.GetYearly()
	result.Monthly = request.GetMonthly()

	return &result, nil
}

func (repo *reportingRepositoryImpl) GetBudgetPerformance(ctx context.Context, request *pb.GetBudgetPerformanceReq) (*pb.GetBudgetPerformanceResp, error) {
	pipeline := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.D{
				{Key: "user_id", Value: request.UserId},
				{Key: "period", Value: "MONTHLY"},
				{Key: "start_date", Value: bson.D{
					{Key: "$gte", Value: time.Date(int(request.Year), time.Month(request.Month), 1, 0, 0, 0, 0, time.UTC)},
					{Key: "$lt", Value: time.Date(int(request.Year), time.Month(request.Month)+1, 1, 0, 0, 0, 0, time.UTC)},
				}},
			},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "transactions"},
				{Key: "localField", Value: "category_id"},
				{Key: "foreignField", Value: "category_id"},
				{Key: "as", Value: "transactions"},
			},
		}},
		bson.D{{
			Key: "$unwind", Value: "$transactions",
		}},
		bson.D{{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$category_id"},
				{Key: "target", Value: bson.D{{Key: "$first", Value: "$amount"}}},
				{Key: "actual", Value: bson.D{{Key: "$sum", Value: "$transactions.amount"}}},
			},
		}},
		bson.D{{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "category_id", Value: "$_id"},
				{Key: "target", Value: 1},
				{Key: "actual", Value: 1},
				{Key: "progress", Value: bson.D{{Key: "$subtract", Value: bson.A{"$target", "$actual"}}}},
			},
		}},
	}

	cursor, err := repo.db.Collection("budgets").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var results []*pb.BudgetPerformance
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &pb.GetBudgetPerformanceResp{
		UserId:                request.UserId,
		Year:                  request.Year,
		Month:                 request.Month,
		BudgetPerformanceList: results,
	}, nil
}

func (repo *reportingRepositoryImpl) GetGoalsProgress(ctx context.Context, request *pb.GetGoalProgressReq) (*pb.GetGoalProgressResp, error) {
	pipeline := mongo.Pipeline{
		bson.D{{
			Key: "$match", Value: bson.D{
				{Key: "user_id", Value: request.UserId},
				{Key: "status", Value: "in_progress"},
			},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "accounts"},
				{Key: "localField", Value: "account_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "account_details"},
			},
		}},
		bson.D{{
			Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$account_details"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		}},
		bson.D{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "transactions"},
				{Key: "localField", Value: "account_id"},
				{Key: "foreignField", Value: "account_id"},
				{Key: "as", Value: "related_transactions"},
			},
		}},
		bson.D{{
			Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$related_transactions"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		}},
		bson.D{{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"},
				{Key: "name", Value: bson.D{{Key: "$first", Value: "$name"}}},
				{Key: "target_amount", Value: bson.D{{Key: "$first", Value: "$target_amount"}}},
				{Key: "current_amount", Value: bson.D{{Key: "$sum", Value: "$related_transactions.amount"}}},
				{Key: "progress", Value: bson.D{{Key: "$multiply", Value: bson.A{bson.D{{Key: "$divide", Value: bson.A{"$current_amount", "$target_amount"}}}, 100}}}},
			},
		}},
	}

	cursor, err := repo.db.Collection("goals").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []*pb.GoalProgress
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return &pb.GetGoalProgressResp{
		GoalProgress: results,
	}, nil
}
