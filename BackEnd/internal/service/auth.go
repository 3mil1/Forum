package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/internal/storage"
	"forum/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	service *Service
	storage *storage.Storage
}

func (s *Session) NewSession(login, pwd string) (string, error) {
	u, err := s.storage.User().FindByLogin(login)
	if err != nil {
		return "", err
	}
	if !u.ComparePassword(u.Password, pwd) {
		return "", appError.InvalidArgumentError(nil, "password is incorrect")
	}
	sessionID := s.GenerateCookieCode()
	if err := s.storage.Auth().CreateSession(u.ID, sessionID); err != nil {
		return "", err
	}
	return sessionID, nil
}

func (s *Session) GenerateCookieCode() string {
	h := hmac.New(sha256.New, []byte("Hello tere privet"))
	newID := uuid.NewV4().Bytes()
	h.Write(newID)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Session) End(userID string) error {
	if err := s.storage.Auth().Delete(userID); err != nil {
		return err
	}
	return nil
}

func (s *Session) Check(cookie, login string) (*models.User, error) {
	user, err := s.storage.Auth().Check(cookie)

	if err != nil {
		logger.InfoLogger.Println(err)
		return nil, err
	}

	if user.Login != login {
		return nil, appError.ForbiddenError
	}

	return user, nil
}
