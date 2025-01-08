package base62

import (
	"fmt"
	"math/big"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	zero = big.NewInt(0)
	base = big.NewInt(62)
)

type encoder struct{}

func New() *encoder {
	return &encoder{}
}

func (e *encoder) Encode(in []byte, length int) (string, error) {
	res := e.encode(in)
	if len(res) < length {
		return "", fmt.Errorf("encoded string < %d", length)
	}
	return res[:length], nil
}

func (e *encoder) encode(in []byte) string {
	num := new(big.Int).SetBytes(in)

	if num.Cmp(zero) == 0 {
		return string(alphabet[0])
	}

	var result []byte
	for num.Cmp(zero) > 0 {
		remainder := new(big.Int)
		num.DivMod(num, base, remainder)
		result = append(result, alphabet[remainder.Int64()])
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
