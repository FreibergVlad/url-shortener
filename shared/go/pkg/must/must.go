package must

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Return[T any](r T, err error) T { //nolint:ireturn
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

func MarshalJSON(v any) []byte {
	return Return(json.Marshal(v))
}

func UnmarshalJSON(data []byte, val any) {
	Do(json.Unmarshal(data, val))
}

func MarshallBSON(v any) []byte {
	return Return(bson.Marshal(v))
}

func UnmarshallBSON(data []byte, val any) {
	Do(bson.Unmarshal(data, val))
}
