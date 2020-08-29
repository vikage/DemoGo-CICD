package authenticate

import (
	"errors"
	"fmt"
	"go-cicd/app/database"
	"go-cicd/app/domain/entity"
	"go-cicd/app/domain/model"
	"go-cicd/app/domain/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ErrTokenInvalid token is invalid error
var ErrTokenInvalid = fmt.Errorf("authenticate: Token invalid")

const encryptKey = "3WkrXalXN4"

// TokenGenerator interface for token generator
type TokenGenerator interface {
	GenTokenForUser(user *entity.User) (string, error)
	GenTokenForUserWithExpireTime(user *entity.User, expirationTime time.Time) (string, error)
}

// TokenDecoder interface for token decoder
type TokenDecoder interface {
	UserFromToken(tokenString string) (*entity.User, error)
}

type tokenGeneratorImpl struct {
}

type tokenDecoderImpl struct {
}

// NewTokenGenerator create token generator
func NewTokenGenerator() TokenGenerator {
	return &tokenGeneratorImpl{}
}

// NewTokenDecoder create token decoder
func NewTokenDecoder() TokenDecoder {
	return &tokenDecoderImpl{}
}

// GenTokenForUser gen jwt for user
func (generator *tokenGeneratorImpl) GenTokenForUser(user *entity.User) (string, error) {
	expirationTime := time.Now().Add(1440 * time.Hour)
	return generator.GenTokenForUserWithExpireTime(user, expirationTime)
}

// GenTokenForUserWithExpireTime gen jwt for user with expire time
func (generator *tokenGeneratorImpl) GenTokenForUserWithExpireTime(user *entity.User, expirationTime time.Time) (string, error) {
	claims := &model.TokenClaims{
		UserID:     user.ID,
		AccountKey: user.AccountKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtKey = []byte(encryptKey)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// UserIDFromToken parse token string and return user id
func (decoder *tokenDecoderImpl) UserFromToken(tokenString string) (*entity.User, error) {
	var jwtKey = []byte(encryptKey)
	claims := &model.TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid == false {
		return nil, ErrTokenInvalid
	}

	dbClient := database.ResolveDatabaseClient()
	userRepo := repository.ResolveUserRepo(dbClient)

	user, err := userRepo.FindUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("User not found")
	}

	if user.AccountKey == claims.AccountKey {
		return user, nil
	}

	return nil, errors.New("Not match account key")
}
