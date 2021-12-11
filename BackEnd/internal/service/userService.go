package service

import (
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/internal/storage"
	"net/mail"
	"regexp"
)

type UserService struct {
	storage *storage.Storage
}

func (u *UserService) Create(user *models.User) (*models.User, error) {
	user.GenerateID()
	if err := u.validate(user); err != nil {
		return nil, err
	}
	user.HashPassword()
	addedUser, err := u.storage.User().Create(user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return addedUser, nil
}

func (u *UserService) SelectAll() ([]models.User, error) {
	users, err := u.storage.User().SelectAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

var (
	minLenPwd = 8
	isCorrect = regexp.MustCompile(`[0-9a-zA-Z]{3,50}$`).MatchString
)

func (u *UserService) validate(user *models.User) error {
	if err := u.validateLogin(user.Login); err != nil {
		return err
	}
	if err := u.validateEmail(user.Email); err != nil {
		return err
	}
	if err := u.validatePwd(user.Password); err != nil {
		return err
	}
	return nil
}

func (u *UserService) validateLogin(login string) error {
	if login == "" {
		return appError.InvalidArgumentError(nil, "login is empty")
	}
	if !isCorrect(login) {
		return appError.InvalidArgumentError(nil, "login format is invalid")
	}
	return nil
}

func (u *UserService) validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return appError.InvalidArgumentError(nil, "email format is invalid")
	}
	return nil
}

func (u *UserService) validatePwd(pwd string) error {
	if pwd == "" {
		return appError.InvalidArgumentError(nil, "password cannot be empty")
	}
	if len(pwd) < minLenPwd {
		return appError.InvalidArgumentError(nil, "password is too short")
	}
	return nil
}

func (u *UserService) FindById(id string) (*models.User, error) {
	user, err := u.storage.User().FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
