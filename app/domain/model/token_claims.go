package model

import "github.com/dgrijalva/jwt-go"

// TokenClaims claims field for fwt token
type TokenClaims struct {
	UserID     string `json:"user_id"`
	AccountKey string `json:"account_key"`
	jwt.StandardClaims
}
