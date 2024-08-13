package mongodb

import (
	pb "budgeting-service/generated/budgeting"
	"budgeting-service/models"
	"budgeting-service/pkg/enums"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *pb.CreateCategoryReq) (*pb.CreateCategoryResp, error)
	UpdateCategory(ctx context.Context, category *pb.UpdateCategoryReq) (*pb.UpdateCategoryResp, error)
	DeleteCategory(ctx context.Context, request *pb.DeleteCategoryReq) (*pb.DeleteCategoryResp, error)
	GetCategory(ctx context.Context, request *pb.GetCategoryReq) (*pb.GetCategoryResp, error)
	GetCategoriesList(ctx context.Context, request *pb.GetCategoriesReq) (*pb.GetCategoriesResp, error)
}

type categoryRepositoryImpl struct {
	coll *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) CategoryRepository {
	return &categoryRepositoryImpl{coll: db.Collection("categories")}
}

func (ropo *categoryRepositoryImpl) CreateCategory(ctx context.Context, category *pb.CreateCategoryReq) (*pb.CreateCategoryResp, error) {
	_, err := ropo.coll.InsertOne(ctx, bson.D{
		{Key: "_id", Value: uuid.NewString()},
		{Key: "user_id", Value: category.UserId},
		{Key: "name", Value: category.Name},
		{Key: "type", Value: category.Type.String()},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
		{Key: "deleted_at", Value: nil},
	})
	if err != nil {
		return &pb.CreateCategoryResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}

	return &pb.CreateCategoryResp{
		Status:  "success",
		Message: "Category created successfully",
	}, nil
}

func (repo *categoryRepositoryImpl) UpdateCategory(ctx context.Context, request *pb.UpdateCategoryReq) (*pb.UpdateCategoryResp, error) {
	updated := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: request.Name},
			{Key: "type", Value: request.Type.String()},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "deleted_at", Value: nil},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, updated)
	if err != nil {
		return &pb.UpdateCategoryResp{
			Status:  "error",
			Message: "Error updating category: " + err.Error(),
		}, err
	}

	if res.ModifiedCount == 0 {
		return &pb.UpdateCategoryResp{
			Status:  "error",
			Message: "Category not found",
		}, nil
	}

	return &pb.UpdateCategoryResp{
		Status:  "success",
		Message: "Category updated successfully",
	}, nil
}

func (repo *categoryRepositoryImpl) DeleteCategory(ctx context.Context, request *pb.DeleteCategoryReq) (*pb.DeleteCategoryResp, error) {
	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "user_id", Value: request.UserId},
		{Key: "deleted_at", Value: nil},
	}

	res, err := repo.coll.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: time.Now()},
		}},
	})
	if err != nil {
		return &pb.DeleteCategoryResp{
			Status:  "error",
			Message: err.Error(),
		}, err
	}

	if res.ModifiedCount == 0 {
		return &pb.DeleteCategoryResp{
			Status:  "error",
			Message: "Category not found",
		}, nil
	}

	return &pb.DeleteCategoryResp{
		Status:  "success",
		Message: "Category deleted successfully",
	}, nil
}

func (repo *categoryRepositoryImpl) GetCategory(ctx context.Context, request *pb.GetCategoryReq) (*pb.GetCategoryResp, error) {
	filter := bson.D{
		{Key: "_id", Value: request.Id},
		{Key: "deleted_at", Value: nil},
	}

	var category models.GetCategory
	err := repo.coll.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	categoryType, err := enums.CategoryTypeParse(category.Type)
	if err != nil {
		return nil, err
	}

	return &pb.GetCategoryResp{
		Id:     category.ID,
		UserId: category.UserId,
		Name:   category.Name,
		Type:   *categoryType,
	}, nil
}

func (repo *categoryRepositoryImpl) GetCategoriesList(ctx context.Context, request *pb.GetCategoriesReq) (*pb.GetCategoriesResp, error) {
	pipeline := createCategoryFilters(request)
	totalCount, err := repo.coll.CountDocuments(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (request.Page - 1) * request.Limit}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: request.Limit}})

	cursor, err := repo.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var categories []*pb.Category
	for cursor.Next(ctx) {
		var category models.GetCategory
		err := cursor.Decode(&category)
		if err != nil {
			return nil, err
		}

		categoryType, err := enums.CategoryTypeParse(category.Type)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &pb.Category{
			Id:     category.ID,
			UserId: category.UserId,
			Name:   category.Name,
			Type:   *categoryType,
		})
	}
	return &pb.GetCategoriesResp{
		Categories: categories,
		TotalCount: totalCount,
		Limit:      request.Limit,
	}, nil
}

func createCategoryFilters(request *pb.GetCategoriesReq) mongo.Pipeline {
	pipeline := mongo.Pipeline{}

	if request.UserId != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: request.UserId},
			}},
		})
	}
	if request.Name != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "name", Value: bson.D{
					{Key: "$regex", Value: "*." + request.Name + ".*"},
					{Key: "$options", Value: "i"},
				}},
			}},
		})
	}
	// if request.Type!= "" {
	// 	pipeline = append(pipeline, bson.D{
	//         {Key: "$match", Value: bson.D{
	//             {Key: "type", Value: request.Type},
	//         }},
	//     })
	// }

	pipeline = append(pipeline, bson.D{
		{Key: "deleted_at", Value: nil},
	})
	return pipeline
}
