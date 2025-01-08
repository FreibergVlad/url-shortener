package must

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Return[T interface{}](r T, err error) T {
	if err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
	return r
}

func Do(err error) {
	if err != nil {
		log.Panic().Err(err).Msg(err.Error())
	}
}

func MarshalJSON(v interface{}) []byte {
	return Return(json.Marshal(v))
}

func UnmarshalJSON(data []byte, val interface{}) {
	Do(json.Unmarshal(data, val))
}

func MarshallBSON(v interface{}) []byte {
	return Return(bson.Marshal(v))
}

func UnmarshallBSON(data []byte, val interface{}) {
	Do(bson.Unmarshal(data, val))
}
