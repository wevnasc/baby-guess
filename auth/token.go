package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type credentials struct {
	email    string
	password string
}

func basicAuth(token string) (*credentials, error) {
	typeToken := strings.Split(token, " ")

	if typeToken[0] != "Basic" {
		return nil, errors.New("invalid authorization type")
	}

	decoded, err := base64.StdEncoding.DecodeString(typeToken[1])

	if err != nil {
		return nil, fmt.Errorf("error to decode token %v", err)
	}

	values := strings.Split(string(decoded), ":")

	if len(values) != 2 {
		return nil, errors.New("token bad formated")
	}

	return &credentials{
		email:    values[0],
		password: values[1],
	}, nil
}

func authToken(account account, secret string, duration time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    account.id.String(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	})

	return claims.SignedString([]byte(secret))
}
