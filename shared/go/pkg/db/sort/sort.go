package sort

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrInvalidOrder = errors.New("'order' should be either '1' (asc) or '-1' (desc)")

func validateOrder(order int) error {
	if order != 1 && order != -1 {
		return ErrInvalidOrder
	}
	return nil
}

type Sort struct {
	Field string
	Order int // 1 for 'asc', -1 for 'desc'
}

func (s Sort) MarshalBSON() ([]byte, error) {
	if err := validateOrder(s.Order); err != nil {
		return nil, err
	}
	return bson.Marshal(bson.D{{Key: s.Field, Value: s.Order}})
}

func (s Sort) MarshalJSON() ([]byte, error) {
	if err := validateOrder(s.Order); err != nil {
		return nil, err
	}
	return json.Marshal(map[string]int{s.Field: s.Order})
}
