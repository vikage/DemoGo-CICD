package entity

import "time"

// User user document structure
type User struct {
	ID          string    `bson:"_id"`
	FirstName   string    `bson:"first_name"`
	LastName    string    `bson:"last_name"`
	Email       string    `bson:"email"`
	Password    string    `bson:"password,omitempty"`
	FacebookID  string    `bson:"facebook_id"`
	GoogleID    string    `bson:"google_id"`
	AppleID     string    `bson:"apple_id"`
	AccountType int       `bson:"account_type"`
	PhoneNumber string    `bson:"phone_number"`
	Avatar      string    `bson:"avatar,omitempty"`
	CreateAt    time.Time `bson:"create_at"`
	Status      int       `bson:"status"`
	AccountKey  string    `bson:"account_key"`
}
