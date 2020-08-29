package model

import (
	"go-cicd/app/domain/entity"
	"time"
)

// User fields of User
// swagger:model
type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar,omitempty"`
	CreateAt    time.Time `json:"create_at"`
	Status      int       `json:"status"`
	AccountType int       `json:"account_type"`
	PhoneNumber string    `json:"phone_number"`
}

// NewUserFromEntity create user from entity
func NewUserFromEntity(entity *entity.User) *User {
	return &User{
		ID:          entity.ID,
		FirstName:   entity.FirstName,
		LastName:    entity.LastName,
		Email:       entity.Email,
		CreateAt:    entity.CreateAt.Truncate(time.Second),
		Status:      entity.Status,
		AccountType: entity.AccountType,
		PhoneNumber: entity.PhoneNumber,
		Avatar:      entity.Avatar,
	}
}
