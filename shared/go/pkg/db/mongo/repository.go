package mongo

import (
	"context"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/filters"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type repository[T any] struct {
	collection *mongo.Collection
}

func NewRepository[T any](collection *mongo.Collection) *repository[T] {
	return &repository[T]{collection: collection}
}

func (r *repository[T]) Create(ctx context.Context, resource *T) error {
	_, err := r.collection.InsertOne(ctx, resource)
	if mongo.IsDuplicateKeyError(err) {
		return errors.ErrDuplicateResource
	}
	return err
}

func (r *repository[T]) FindOne(ctx context.Context, filter filters.Filter) (*T, error) {
	var resource T
	err := r.collection.FindOne(ctx, filter).Decode(&resource)
	if err == mongo.ErrNoDocuments {
		return nil, errors.ErrResourceNotFound
	}
	return &resource, err
}

func (r *repository[T]) FindMany(ctx context.Context, query db.PaginatedQuery) (*db.PaginatedResult[*T], error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: query.Filter}},
		{{Key: "$facet", Value: bson.M{
			"documents": []bson.M{
				{"$skip": (query.PageNum - 1) * query.PageSize},
				{"$limit": query.PageSize},
				{"$sort": query.Sort},
			},
			"total": []bson.M{
				{"$count": "count"},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	documents := result[0]["documents"].(bson.A)
	totalCount := int32(0)
	if total := result[0]["total"].(bson.A); len(total) > 0 {
		totalCount = total[0].(bson.D)[0].Value.(int32)
	}

	resources := make([]*T, len(documents))
	for i, doc := range documents {
		var resource T
		bytes, err := bson.Marshal(doc)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(bytes, &resource)
		if err != nil {
			return nil, err
		}
		resources[i] = &resource
	}

	return &db.PaginatedResult[*T]{
		Data:  resources,
		Total: totalCount,
	}, nil
}

func (r *repository[T]) ReplaceOne(ctx context.Context, filter filters.Filter, resource *T) error {
	_, err := r.collection.ReplaceOne(ctx, filter, resource)
	if err == mongo.ErrNoDocuments {
		return errors.ErrResourceNotFound
	}
	return err
}

func (r *repository[T]) UpdateOne(ctx context.Context, filter filters.Filter, update map[string]any) (*T, error) {
	var resource T
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options).Decode(&resource)
	if err == mongo.ErrNoDocuments {
		return nil, errors.ErrResourceNotFound
	}
	return &resource, err
}
