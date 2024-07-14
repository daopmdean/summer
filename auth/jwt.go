package auth

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func GenToken(claims SummerClaim, signedKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(signedKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseToken(token string, signedKey string) (*SummerClaim, error) {
	var claims SummerClaim

	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("wrong signing alg")
		}
		return []byte(signedKey), nil
	})

	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func ExtractTokenFromHeader(bearer string) (string, error) {
	authToken := strings.Split(bearer, " ")
	if len(authToken) != 2 {
		return "", errors.New("invalid bearer token")
	}

	if authToken[0] != "Bearer" {
		return "", errors.New("invalid bearer token")
	}

	return authToken[1], nil
}
