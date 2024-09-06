package services

import (
	"errors"
	"integrand/persistence"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GetUsers(userId int) ([]persistence.User, error) {
	err := isUserAuthorized(userId, persistence.READ_USER)
	if err != nil {
		return []persistence.User{}, err
	}
	return persistence.DATASTORE.GetUsers()
}

func GetUser(callerUserId int, userIDToGet int) (persistence.User, error) {
	err := isUserAuthorized(callerUserId, persistence.READ_USER)
	if err != nil {
		return persistence.User{}, err
	}
	return persistence.DATASTORE.GetUserByID(userIDToGet)
}

func UpdatePassword(callerUserId int, userIdToUpdate int, oldPasswordPlain string, newPasswordPlain string) (persistence.User, error) {
	var user persistence.User

	// If caller is the same as the user to update we can proceed
	if callerUserId != userIdToUpdate {
		// Check if the user has write priviledges
		err := isUserAuthorized(callerUserId, persistence.WRITE_USER)
		if err != nil {
			return persistence.User{}, err
		}
	}

	user, err := persistence.DATASTORE.GetUserByID(userIdToUpdate)
	if err != nil {
		slog.Error(err.Error())
		return persistence.User{}, err
	}
	if checkPasswordHash(string(oldPasswordPlain), user.Password) {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(newPasswordPlain), bcrypt.DefaultCost)
		if err != nil {
			return user, err
		}
		return persistence.DATASTORE.UpdateEmailUser(userIdToUpdate, string(newPassword))
	} else {
		return persistence.User{}, errors.New("password not valid")
	}
}

func CreateUser(userId int, email string, plainPassword string) (persistence.User, error) {
	var user persistence.User
	err := isUserAuthorized(userId, persistence.WRITE_USER)
	if err != nil {
		return persistence.User{}, err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(err.Error())
		return user, err
	}
	user = persistence.User{
		Email:        email,
		Password:     string(password),
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
	}

	id, err := persistence.DATASTORE.CreateUser(user)
	if err != nil {
		return user, err
	}

	user.ID = id
	return user, nil
}

func RemoveUser(callerUserId int, userIDToDelete int) (int, error) {
	err := isUserAuthorized(callerUserId, persistence.WRITE_USER)
	if err != nil {
		return 0, err
	}
	if userIDToDelete == 1 {
		return -1, errors.New("can't delete root user")
	}
	return persistence.DATASTORE.DeleteUserByID(userIDToDelete)
}
