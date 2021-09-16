package token

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const Secret string = "secret"

const (
	Bearer string = "Bearer"
	Basic         = "Basic"
)

type Credentials struct {
	Email    string
	Password string
}

func extractToken(token string, tokenType string) (string, error) {
	typeToken := strings.Split(token, " ")

	if len(typeToken) != 2 {
		return "", errors.New("token bad formated")
	}

	if typeToken[0] != tokenType {
		return "", errors.New("invalid authorization type")
	}

	return typeToken[1], nil
}

func BasicAuth(token string) (*Credentials, error) {
	basic, err := extractToken(token, Basic)

	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(basic)

	if err != nil {
		return nil, fmt.Errorf("error to decode token %v", err)
	}

	values := strings.Split(string(decoded), ":")

	if len(values) != 2 {
		return nil, errors.New("token bad formated")
	}

	return &Credentials{
		Email:    values[0],
		Password: values[1],
	}, nil
}

func Auth(token string, secret string) (issuerID uuid.UUID, err error) {
	bearer, err := extractToken(token, Bearer)

	jwtToken, err := jwt.ParseWithClaims(bearer, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return issuerID, fmt.Errorf("error to decode token %v", err)
	}

	claims := jwtToken.Claims.(*jwt.StandardClaims)

	issuerID, err = uuid.Parse(claims.Issuer)

	if err != nil {
		return issuerID, fmt.Errorf("error to parse uuid %v", err)
	}

	return issuerID, nil

}

func NewAuth(issuer uuid.UUID, secret string, duration time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer.String(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	})

	return claims.SignedString([]byte(secret))
}
