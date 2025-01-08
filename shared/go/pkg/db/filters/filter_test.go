package filters_test

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/filters"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/stretchr/testify/assert"
)

func TestEqualsFilter(t *testing.T) {
	t.Parallel()

	// name = "John"
	filter := filters.Equals[string]{Field: "name", Value: "John"}
	expectedBson := bson.D{{Key: "name", Value: "John"}}

	var actualBson bson.D
	must.UnmarshallBSON(must.MarshallBSON(filter), &actualBson)

	assert.Equal(t, expectedBson, actualBson)
	assert.Equal(t, expectedBson.String(), filter.String())
}

func TestInRangeFilter(t *testing.T) {
	t.Parallel()

	min := must.Return(time.Parse(time.DateOnly, "2024-10-03"))
	max := must.Return(time.Parse(time.DateOnly, "2024-09-03"))
	// date <= 2024-10-03 AND date >= 2024-09-03
	filter := filters.InRange[time.Time]{Field: "date", Min: min, Max: max}
	expectedBson := bson.D{
		{Key: "date", Value: bson.D{
			{Key: "$gte", Value: bson.DateTime(min.UnixMilli())},
			{Key: "$lte", Value: bson.DateTime(max.UnixMilli())},
		}},
	}

	var actualBson bson.D
	must.UnmarshallBSON(must.MarshallBSON(filter), &actualBson)

	assert.Equal(t, expectedBson, actualBson)
	assert.Equal(t, expectedBson.String(), filter.String())
}

func TestInFilter(t *testing.T) {
	t.Parallel()

	// _id IN ["1", "2", "3"]
	filter := filters.In[string]{Field: "_id", Values: []string{"1", "2", "3"}}
	expectedBson := bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "$in", Value: bson.A{"1", "2", "3"}},
		}},
	}

	var actualBson bson.D
	must.UnmarshallBSON(must.MarshallBSON(filter), &actualBson)

	assert.Equal(t, expectedBson, actualBson)
	assert.Equal(t, expectedBson.String(), filter.String())
}

func TestAndFilter(t *testing.T) {
	t.Parallel()

	min, max := int32(18), int32(30)
	// name = "John" AND age IN [18, 30]
	filter := filters.And{
		Filters: []filters.Filter{
			filters.Equals[string]{Field: "name", Value: "John"},
			filters.InRange[int32]{Field: "age", Min: min, Max: max},
		},
	}
	expectedBson := bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "name", Value: "John"}},
			bson.D{{Key: "age", Value: bson.D{
				{Key: "$gte", Value: min},
				{Key: "$lte", Value: max},
			}}},
		}},
	}

	var actualBson bson.D
	must.UnmarshallBSON(must.MarshallBSON(filter), &actualBson)

	assert.Equal(t, expectedBson, actualBson)
	assert.Equal(t, expectedBson.String(), filter.String())
}

func TestOrFilter(t *testing.T) {
	t.Parallel()

	// name = "John" OR name = "Jane"
	filter := filters.Or{
		Filters: []filters.Filter{
			filters.Equals[string]{Field: "name", Value: "John"},
			filters.Equals[string]{Field: "name", Value: "Jane"},
		},
	}
	expectedBson := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: "John"}},
			bson.D{{Key: "name", Value: "Jane"}},
		}},
	}

	var actualBson bson.D
	must.UnmarshallBSON(must.MarshallBSON(filter), &actualBson)

	assert.Equal(t, expectedBson, actualBson)
	assert.Equal(t, expectedBson.String(), filter.String())
}

func TestNorFilter(t *testing.T) {
	t.Parallel()

	// name != "John"
	filter := filters.Nor{
		Filters: []filters.Filter{filters.Equals[string]{Field: "name", Value: "John"}},
	}
	expectedBson := bson.D{{Key: "$nor", Value: bson.A{bson.D{{Key: "name", Value: "John"}}}}}

	var actualBson bson.D
	must.UnmarshallBSON(must.MarshallBSON(filter), &actualBson)

	assert.Equal(t, expectedBson, actualBson)
	assert.Equal(t, expectedBson.String(), filter.String())
}

func TestComplexFilter(t *testing.T) {
	t.Parallel()

	min, max := int32(18), int32(30)
	// (name = "John" OR name = "Jane") AND NOT age in [18, 30]
	filter := filters.And{
		Filters: []filters.Filter{
			filters.Or{
				Filters: []filters.Filter{
					filters.Equals[string]{Field: "name", Value: "John"},
					filters.Equals[string]{Field: "name", Value: "Jane"},
				},
			},
			filters.Nor{
				Filters: []filters.Filter{filters.InRange[int32]{Field: "age", Min: min, Max: max}},
			},
		},
	}
	expectedBson := bson.D{
		{
			Key: "$and",
			Value: bson.A{
				bson.D{{
					Key: "$or",
					Value: bson.A{
						bson.D{{Key: "name", Value: "John"}},
						bson.D{{Key: "name", Value: "Jane"}},
					}},
				},
				bson.D{{
					Key: "$nor",
					Value: bson.A{
						bson.D{{
							Key: "age",
							Value: bson.D{
								{Key: "$gte", Value: min},
								{Key: "$lte", Value: max},
							},
						}},
					},
				}},
			},
		},
	}

	var actualBson bson.D
	must.UnmarshallBSON(must.MarshallBSON(filter), &actualBson)

	assert.Equal(t, expectedBson, actualBson)
	assert.Equal(t, expectedBson.String(), filter.String())
}
