package sort_test

import (
	"encoding/json"
	"testing"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/db/sort"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestSort(t *testing.T) {
	sort := sort.Sort{Field: "field", Order: 1}

	var actualBson bson.M
	must.UnmarshallBSON(must.MarshallBSON(sort), &actualBson)

	order := int32(1)
	expectedBson := bson.M{"field": order}
	assert.Equal(t, expectedBson, actualBson)

	expectedJson := `{"field": 1}`
	assert.JSONEq(t, expectedJson, string(must.MarshalJSON(sort)))
}

func TestSortWithInvalidOrder(t *testing.T) {
	sorter := sort.Sort{Field: "field", Order: 0}

	_, err := bson.Marshal(sorter)
	assert.ErrorIs(t, err, sort.ErrInvalidOrder)

	_, err = json.Marshal(sorter)
	assert.ErrorIs(t, err, sort.ErrInvalidOrder)
}
