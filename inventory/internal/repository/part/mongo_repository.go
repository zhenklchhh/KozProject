package part

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zhenklchhh/KozProject/inventory/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	collection := db.Collection("parts")
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "manufacturer.country", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}
	return &MongoRepository{
		collection: db.Collection("parts"),
	}, nil
}

func (r *MongoRepository) GetPart(ctx context.Context, id uuid.UUID) (*model.Part, error) {
	var part model.Part
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &model.Part{}, model.ErrNotFound
		}
		return &model.Part{}, err
	}
	return &part, nil
}

func (r *MongoRepository) ListParts(ctx context.Context, partFilter *model.PartFilter) ([]*model.Part, error) {
	filter := bson.M{}

	if partFilter == nil {
		return r.find(ctx, filter)
	}

	if len(partFilter.Uuids) > 0 {
		filter["_id"] = bson.M{"$in": partFilter.Uuids}
	}
	if len(partFilter.Names) > 0 {
		filter["name"] = bson.M{"$in": partFilter.Names}
	}
	if len(partFilter.Tags) > 0 {
		filter["tags"] = bson.M{"$in": partFilter.Tags}
	}
	if len(partFilter.Categories) > 0 {
		filter["category"] = bson.M{"$in": partFilter.Categories}
	}
	if len(partFilter.ManufacturerCountries) > 0 {
		filter["manufacturer.country"] = bson.M{"$in": partFilter.ManufacturerCountries}
	}
	return r.find(ctx, filter)
}

func (r *MongoRepository) find(ctx context.Context, filter bson.M) ([]*model.Part, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var parts []*model.Part
	if err := cursor.All(ctx, &parts); err != nil {
		return nil, err
	}
	return parts, nil
}
