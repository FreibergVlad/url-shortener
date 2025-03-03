package shorturls

import (
	"context"
	"errors"

	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrAlreadyExists = errors.New("resource already exists")
	ErrNotFound      = errors.New("resource not found")
)

type Repository interface {
	Create(ctx context.Context, shorturl *schema.ShortURL) error
	GetByGlobalKey(ctx context.Context, key schema.ShortURLGlobalKey) (*schema.ShortURL, error)
	GetByOrganizationIDAndGlobalKey(
		ctx context.Context,
		organizationID string,
		key schema.ShortURLGlobalKey,
	) (*schema.ShortURL, error)
	GetByOrganizationKeyOrGlobalKey(
		ctx context.Context,
		orgKey schema.ShortURLOrganizationKey,
		globalKey schema.ShortURLGlobalKey,
	) (*schema.ShortURL, error)
	ListByOrganizationID(
		ctx context.Context,
		organizationID string,
		pageNum,
		pageSize int,
	) (*PaginatedResult, error)
	ReplaceByOrganizationIDAndGlobalKey(
		ctx context.Context,
		organizationID string,
		key schema.ShortURLGlobalKey,
		replaceWith *schema.ShortURL,
	) error
	DeleteByOrganizationIDAndGlobalKey(
		ctx context.Context,
		organizationID string,
		key schema.ShortURLGlobalKey,
	) (*schema.ShortURL, error)
}

type PaginatedResult struct {
	Data  []*schema.ShortURL
	Total int32
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{collection: collection}
}

func (r *MongoRepository) Create(ctx context.Context, shorturl *schema.ShortURL) error {
	_, err := r.collection.InsertOne(ctx, shorturl)
	if mongo.IsDuplicateKeyError(err) {
		return ErrAlreadyExists
	}
	return err
}

func (r *MongoRepository) GetByGlobalKey(ctx context.Context, key schema.ShortURLGlobalKey) (*schema.ShortURL, error) {
	var shorturl schema.ShortURL
	err := r.collection.FindOne(ctx, globalKeyFilter(key)).Decode(&shorturl)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	return &shorturl, err
}

func (r *MongoRepository) GetByOrganizationIDAndGlobalKey(
	ctx context.Context, organizationID string, key schema.ShortURLGlobalKey,
) (*schema.ShortURL, error) {
	var shorturl schema.ShortURL
	err := r.collection.
		FindOne(ctx, organizationIDAndGlobalKeyFilter(organizationID, key)).
		Decode(&shorturl)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	return &shorturl, err
}

func (r *MongoRepository) GetByOrganizationKeyOrGlobalKey(
	ctx context.Context, orgKey schema.ShortURLOrganizationKey, globalKey schema.ShortURLGlobalKey,
) (*schema.ShortURL, error) {
	filter := bson.M{"$or": bson.A{globalKeyFilter(globalKey), organizationKeyFilter(orgKey)}}
	var shorturl schema.ShortURL
	err := r.collection.FindOne(ctx, filter).Decode(&shorturl)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	return &shorturl, err
}

func (r *MongoRepository) ListByOrganizationID(
	ctx context.Context, organizationID string, pageNum, pageSize int,
) (*PaginatedResult, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"organization_id": organizationID}}},
		{{Key: "$facet", Value: bson.M{
			"documents": []bson.M{
				{"$skip": (pageNum - 1) * pageSize},
				{"$limit": pageSize},
				{"$sort": bson.D{{Key: "created_at", Value: -1}}},
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

	shorturls := make([]*schema.ShortURL, len(documents))
	for docIndex, doc := range documents {
		var shorturl schema.ShortURL
		bytes, err := bson.Marshal(doc)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(bytes, &shorturl)
		if err != nil {
			return nil, err
		}
		shorturls[docIndex] = &shorturl
	}

	return &PaginatedResult{
		Data:  shorturls,
		Total: totalCount,
	}, nil
}

func (r *MongoRepository) ReplaceByOrganizationIDAndGlobalKey(
	ctx context.Context, organizationID string, key schema.ShortURLGlobalKey, replaceWith *schema.ShortURL,
) error {
	var shorturl schema.ShortURL
	err := r.collection.
		FindOneAndReplace(
			ctx,
			organizationIDAndGlobalKeyFilter(organizationID, key),
			replaceWith,
		).
		Decode(&shorturl)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNotFound
	}
	if mongo.IsDuplicateKeyError(err) {
		return ErrAlreadyExists
	}
	return err
}

func (r *MongoRepository) DeleteByOrganizationIDAndGlobalKey(
	ctx context.Context, organizationID string, key schema.ShortURLGlobalKey,
) (*schema.ShortURL, error) {
	var shorturl schema.ShortURL
	err := r.collection.
		FindOneAndDelete(ctx, organizationIDAndGlobalKeyFilter(organizationID, key)).
		Decode(&shorturl)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	return &shorturl, err
}

func globalKeyFilter(key schema.ShortURLGlobalKey) bson.M {
	return bson.M{"$and": bson.A{
		bson.M{"domain": key.Domain},
		bson.M{"alias": key.Alias},
	}}
}

func organizationIDAndGlobalKeyFilter(organizationID string, key schema.ShortURLGlobalKey) bson.M {
	return bson.M{"$and": bson.A{
		bson.M{"domain": key.Domain},
		bson.M{"alias": key.Alias},
		bson.M{"organization_id": organizationID},
	}}
}

func organizationKeyFilter(key schema.ShortURLOrganizationKey) bson.M {
	return bson.M{"$and": bson.A{
		bson.M{"organization_id": key.OrganizationID},
		bson.M{"domain": key.Domain},
		bson.M{"long_url.hash": key.LongURLHash},
	}}
}
