package jwt

import (
	"time"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/golang-jwt/jwt/v5"
)

func IssueForUserId(userId, secret string, issuedAt time.Time, lifetimeSeconds int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(issuedAt.Add(time.Second * time.Duration(lifetimeSeconds))),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(secret))
}

func VerifyAndParseUserId(rawToken, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(
		rawToken,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{"HS512"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}
	return must.Return(token.Claims.GetSubject()), nil
}
