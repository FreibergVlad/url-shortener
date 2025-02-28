package jwt

import (
	"errors"
	"time"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenExpired = errors.New("JWT token is expired")

func IssueForUserID(userID, secret string, issuedAt time.Time, lifetimeSeconds int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(issuedAt.Add(time.Second * time.Duration(lifetimeSeconds))),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(secret))
}

func VerifyAndParseUserID(rawToken, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(
		rawToken,
		&jwt.RegisteredClaims{},
		func(_ *jwt.Token) (interface{}, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{"HS512"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}
		return "", err
	}
	return must.Return(token.Claims.GetSubject()), nil
}
