package model

// User struct in Go
type User struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Friends     []string `json:"friends" db:"friends"`
	IsOnline    bool     `json:"isOnline" db:"isOnline" bson:"isOnline"`
	IsTyping    bool     `json:"isTyping" db:"isTyping" bson:"isTyping"`
	PhoneNumber string   `json:"phoneNumber" db:"phoneNumber" bson:"phoneNumber"`
	CreatedAt   string   `json:"createdAt" db:"createdAt"`
}
