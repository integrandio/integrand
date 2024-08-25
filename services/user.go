package services

import (
	"errors"
	"integrand/persistence"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func GetUsers() ([]persistence.User, error) {
	return persistence.DATASTORE.GetEmailUsers()
}

func GetUser(id int) (persistence.User, error) {
	return persistence.DATASTORE.GetUserByID(id)
}

func UpdatePassword(id int, oldPasswordPlain string, newPasswordPlain string) (persistence.User, error) {
	var user persistence.User
	user, err := persistence.DATASTORE.GetUserByID(id)
	if err != nil {
		slog.Error(err.Error())
		return persistence.User{}, err
	}
	if checkPasswordHash(string(oldPasswordPlain), user.Password) {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(newPasswordPlain), bcrypt.DefaultCost)
		if err != nil {
			return user, err
		}
		return persistence.DATASTORE.UpdateEmailUser(id, string(newPassword))
	} else {
		return persistence.User{}, errors.New("password not valid")
	}
}

func CreateUser(email string, plainPassword string, auth_type persistence.AuthType, socialId string) (persistence.User, error) {
	var user persistence.User
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(err.Error())
		return user, err
	}
	user = persistence.User{
		Email:    email,
		Password: string(password),
		AuthType: auth_type,
		SocialID: socialId,
	}

	id, err := persistence.DATASTORE.CreateEmailUser(user)
	if err != nil {
		return user, err
	}

	user.ID = id
	return user, nil
}

func RemoveUser(id int) (int, error) {
	if id == 1 {
		return -1, errors.New("can't delete root user")
	}
	return persistence.DATASTORE.DeleteEmailUser(id)
}
