package filters

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

func String(filter Filter) string {
	s, err := bson.MarshalExtJSON(filter, true, false)
	if err != nil {
		return ""
	}
	return string(s)
}

type Filter interface {
	bson.Marshaler
	String() string
}

type Equals[T any] struct {
	Field string
	Value T
}

type InRange[T any] struct {
	Field string
	Min   T
	Max   T
}

type In[T any] struct {
	Field  string
	Values []T
}

type And struct {
	Filters []Filter
}

type Or struct {
	Filters []Filter
}

type Nor struct {
	Filters []Filter
}

func (e Equals[T]) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.D{{Key: e.Field, Value: e.Value}})
}

func (e Equals[T]) String() string {
	return String(e)
}

func (r InRange[T]) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.D{{
		Key: r.Field,
		Value: bson.D{
			{Key: "$gte", Value: r.Min},
			{Key: "$lte", Value: r.Max},
		},
	}})
}

func (r InRange[T]) String() string {
	return String(r)
}

func (i In[T]) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.D{{Key: i.Field, Value: bson.D{{Key: "$in", Value: i.Values}}}})
}

func (i In[T]) String() string {
	return String(i)
}

func (a And) MarshalBSON() ([]byte, error) {
	andFilters := make([]bson.D, len(a.Filters))
	for i, filter := range a.Filters {
		filterBSON, err := filter.MarshalBSON()
		if err != nil {
			return nil, err
		}
		var doc bson.D
		if err := bson.Unmarshal(filterBSON, &doc); err != nil {
			return nil, err
		}
		andFilters[i] = doc
	}

	return bson.Marshal(bson.D{{Key: "$and", Value: andFilters}})
}

func (a And) String() string {
	return String(a)
}

func (o Or) MarshalBSON() ([]byte, error) {
	orFilters := make([]bson.D, len(o.Filters))
	for i, filter := range o.Filters {
		filterBSON, err := filter.MarshalBSON()
		if err != nil {
			return nil, err
		}
		var doc bson.D
		if err := bson.Unmarshal(filterBSON, &doc); err != nil {
			return nil, err
		}
		orFilters[i] = doc
	}

	return bson.Marshal(bson.D{{Key: "$or", Value: orFilters}})
}

func (o Or) String() string {
	return String(o)
}

func (n Nor) MarshalBSON() ([]byte, error) {
	norFilters := make([]bson.D, len(n.Filters))
	for i, filter := range n.Filters {
		filterBSON, err := filter.MarshalBSON()
		if err != nil {
			return nil, err
		}
		var doc bson.D
		if err := bson.Unmarshal(filterBSON, &doc); err != nil {
			return nil, err
		}
		norFilters[i] = doc
	}
	return bson.Marshal(bson.D{{Key: "$nor", Value: norFilters}})
}

func (n Nor) String() string {
	return String(n)
}
