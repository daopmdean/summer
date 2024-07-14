package auth

import "github.com/golang-jwt/jwt/v5"

type SummerClaim struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
