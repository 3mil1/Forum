package storage

import (
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/pkg/logger"
	"time"
)

type AuthRepo struct {
	storage *Storage
}

var (
	sessionColumns = "session_key, user_id, expired_at"
)

func (auth *AuthRepo) CreateSession(userID, sessionKey string) error {
	query := fmt.Sprintf("INSERT INTO sessions (%s) VALUES ($1, $2, $3)", sessionColumns)
	t := time.Now().Add(3 * 24 * time.Hour)
	_, err := auth.storage.db.Exec(query, sessionKey, userID, t)

	if err != nil {
		// update session_key if exist
		if err.Error() == "UNIQUE constraint failed: sessions.user_id" {
			query := fmt.Sprintf("UPDATE sessions SET session_key=$1 WHERE user_id=$2")
			t := time.Now().Add(3 * 24 * time.Hour)
			_, err := auth.storage.db.Exec(query, sessionKey, userID, t)
			if err != nil {
				return appError.SystemError(err)
			}
			return nil
		}
		return appError.SystemError(err)
	}
	return nil
}

func (auth *AuthRepo) Delete(userID string) error {
	query := fmt.Sprintf("DELETE FROM sessions WHERE user_id=$1")
	_, err := auth.storage.db.Exec(query, userID)
	if err != nil {
		return appError.InvalidArgumentError(err, "no current session")
	}
	return nil
}

func (auth *AuthRepo) Check(sessionKey string) (*models.User, error) {
	query := fmt.Sprintf("SELECT u.id, u.email, u.login FROM sessions INNER JOIN users u on u.id = sessions.user_id WHERE session_key=$1")
	row := auth.storage.db.QueryRow(query, sessionKey)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Login)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, appError.InvalidArgumentError(err, "no current session")
	}
	return &user, nil
}
