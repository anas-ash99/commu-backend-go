package model

// Chat struct in Go
type Chat struct {
	ID               string   `json:"id" db:"id" bson:"id"`
	Users            []string `json:"users" db:"users" bson:"users"`
	CreatedAt        string   `json:"createdAt" db:"createdAt" bson:"createdAt"`
	UserChattingWith User     `json:"userChattingWith"`
}
