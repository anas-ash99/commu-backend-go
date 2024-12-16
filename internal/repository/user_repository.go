package repository

import (
	"message-broker/internal/db"
	"message-broker/internal/utils"
	"message-broker/model"
)

func CreateUser(user *model.User) error {
	user.CreatedAt = utils.GetCurrentTime()
	err := db.UpsertUser(user)
	return err
}

func GetAllUsers() ([]*model.User, error) {
	users, err := db.FindAllUsers()
	if users == nil {
		users = []*model.User{}
	}
	return users, err
}

func AddFriend(userId string, friendId string) error {
	return db.AddFriend(userId, friendId)
}

func GetFriendsForUser(userId string) ([]*model.User, error) {
	users, err := db.FindFriendsForUser(userId)
	if users == nil {
		users = []*model.User{}
	}
	return users, err
}

func DeleteUser(user *model.User) error {
	return nil
}

func UpdateUser(user *model.User) error {
	return db.UpdateUser(user)
}
