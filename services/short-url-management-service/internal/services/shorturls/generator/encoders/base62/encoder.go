package base62

import (
	"errors"
	"math/big"
)

const (
	Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Base     = 62
)

var (
	bigZero            = big.NewInt(0)
	bigBase            = big.NewInt(Base)
	errInsufficientLen = errors.New("encoded string < requested length")
)

type Encoder struct{}

func New() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(in []byte, length int) (string, error) {
	res := e.encode(in)
	if len(res) < length {
		return "", errInsufficientLen
	}
	return res[:length], nil
}

func (e *Encoder) encode(in []byte) string {
	num := new(big.Int).SetBytes(in)

	if num.Cmp(bigZero) == 0 {
		return string(Alphabet[0])
	}

	var result []byte
	for num.Cmp(bigZero) > 0 {
		remainder := new(big.Int)
		num.DivMod(num, bigBase, remainder)
		result = append(result, Alphabet[remainder.Int64()])
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
