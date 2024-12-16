package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"message-broker/model"
)

var mongoClient *mongo.Client

// ConnectToMongoDB initializes the MongoDB mongoClient and connects to the database.
func ConnectToMongoDB(uri string) error {
	if mongoClient != nil {
		// If the mongoClient is already initialized, no need to reconnect.
		return nil
	}

	// Set the mongoClient options
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	// Connect to MongoDB
	mongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")
	return nil
}

func InsertMessage(msg *model.Message) error {
	collection := mongoClient.Database("commu").Collection("messages")

	filter := bson.M{"id": msg.ID}           // Match document with the same 'id'
	update := bson.M{"$set": msg}            // Update the document fields with new data
	opts := options.Update().SetUpsert(true) // Enable upsert option

	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	fmt.Println("Inserted a document into the message collection")
	return nil
}

// FineMessagesByField fetches all messages for a specific user by userId.
func FineMessagesByField(value string, key string) ([]*model.Message, error) {
	collection := mongoClient.Database("commu").Collection("messages")

	var filter bson.M
	if key == "userId" {
		filter = bson.M{
			"$or": []bson.M{
				{"authorUserId": value},
				{"otherUserId": value},
			},
		}

	} else {
		filter = bson.M{
			"$and": []bson.M{
				{key: value},
			},
		}
	}

	// Create a slice to hold the resulting messages
	var messages []*model.Message

	// Use Find to retrieve documents matching the filter
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		// Log the error and return
		log.Printf("Error finding messages   %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and decode the documents into the messages slice
	for cursor.Next(context.Background()) {
		var msg model.Message
		if err := cursor.Decode(&msg); err != nil {
			log.Printf("Error decoding message: %v", err)
			return nil, err
		}
		// Append the decoded message to the result slice
		messages = append(messages, &msg)
	}

	// Check if any error occurred during iteration
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	// Return the messages
	return messages, nil
}

func InsertChat(chat *model.Chat) error {

	collection := mongoClient.Database("commu").Collection("chats")

	_, err := collection.InsertOne(context.Background(), chat)
	if err != nil {
		// Log the error for debugging and return it
		log.Printf("Error inserting chat: %v", err)
		return fmt.Errorf("failed to insert chat: %w", err)
	}
	return nil
}

func FindChatsForUser(userId string) ([]*model.Chat, error) {
	chatCollection := mongoClient.Database("commu").Collection("chats")
	//userCollection := mongoClient.Database("commu").Collection("users")
	cursor, err := chatCollection.Find(context.Background(), bson.M{"users": userId})

	defer cursor.Close(context.Background())

	if err != nil {
		log.Printf("Error finding chats: %v", err)
		return nil, err
	}

	var chats []*model.Chat
	if err = cursor.All(context.TODO(), &chats); err != nil {
		log.Fatalf("Error decoding chats: %v", err)
	}
	return chats, nil

}

func AddFriend(userId string, friendId string) error {
	collection := mongoClient.Database("commu").Collection("users")
	var user model.User
	filter := bson.M{"id": userId}

	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return err
	}
	user.Friends = append(user.Friends, friendId)
	_, err = collection.ReplaceOne(context.Background(), filter, user)
	if err != nil {
		return err
	}
	return nil
}

func findUsersByIDs(userIDs []string) ([]*model.User, error) {

	collection := mongoClient.Database("commu").Collection("users")
	ctx := context.TODO()

	// Create a filter using the $in operator
	filter := bson.M{"id": bson.M{"$in": userIDs}}

	// Find matching users
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice of User structs
	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func UpsertUser(user *model.User) error {
	collection := mongoClient.Database("commu").Collection("users")

	filter := bson.M{"id": user.ID}          // Match document with the same 'id'
	update := bson.M{"$set": user}           // Update the document fields with new data
	opts := options.Update().SetUpsert(true) // Enable upsert option

	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}

	return nil
}

func FindAllUsers() ([]*model.User, error) {
	collection := mongoClient.Database("commu").Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error finding users: %v", err)
		return nil, err
	}

	var users []*model.User
	if err = cursor.All(context.Background(), &users); err != nil {
		log.Fatalf("Error decoding users: %v", err)
		return nil, err
	}
	return users, nil
}

func FindFriendsForUser(userId string) ([]*model.User, error) {
	collection := mongoClient.Database("commu").Collection("users")

	var user *model.User

	err := collection.FindOne(context.Background(), bson.M{"id": userId}).Decode(&user)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, err
	}
	fmt.Println(user.Friends)
	users, err := findUsersByIDs(user.Friends)

	if err != nil {
		log.Printf("Error finding users: %v", err)
		return nil, err
	}
	return users, nil

}

func UpdateUser(user *model.User) error {
	collection := mongoClient.Database("commu").Collection("users")

	filter := bson.M{"id": user.ID}          // Match document with the same 'id'
	update := bson.M{"$set": user}           // Update the document fields with new data
	opts := options.Update().SetUpsert(true) // Enable upsert option

	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}
	return nil
}
