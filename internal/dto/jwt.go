package dto

import "github.com/golang-jwt/jwt/v5"

type JwtPayLoad struct {
	jwt.RegisteredClaims
	UserResponse
}
