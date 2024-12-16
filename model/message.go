package model

// Message struct in Go
type Message struct {
	ID           string `json:"id" db:"id" bson:"id"`
	AuthorUserId string `json:"authorUserId" db:"authorUserId" bson:"authorUserId"`
	OtherUserId  string `json:"otherUserId" db:"otherUserId" bson:"otherUserId"`
	ChatID       string `json:"chatId" db:"chatId" bson:"chatId"`
	Content      string `json:"content" db:"content" bson:"content"`
	CreatedAt    string `json:"createdAt" db:"createdAt" bson:"createdAt"`
	DeliveredAt  string `json:"deliveredAt,omitempty" db:"deliveredAt,omitempty" bson:"deliveredAt,omitempty"`
	SeenAt       string `json:"seenAt,omitempty" db:"seenAt,omitempty" bson:"seenAt,omitempty"`
}
