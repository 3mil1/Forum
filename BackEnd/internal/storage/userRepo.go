package storage

import (
	"errors"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

type UserRepo struct {
	storage *Storage
}

var (
	userColumns = "id, email, login, password"
)

// Create User in db
func (ur *UserRepo) Create(u *models.User) (*models.User, error) {
	query := fmt.Sprintf("INSERT INTO users (%s) VALUES ($1, $2, $3, $4)", userColumns)
	if _, err := ur.storage.db.Exec(query, u.ID, u.Email, u.Login, u.Password); err != nil {
		fmt.Println(err)
		var sErr sqlite3.Error
		if errors.As(err, &sErr) {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
				return nil, appError.InvalidArgumentError(nil, "user with this email already exists")
			}
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.login") {
				return nil, appError.InvalidArgumentError(nil, "user with this login already exists")
			}
		}
		return nil, appError.SystemError(err)
	}
	return u, nil
}

func (ur *UserRepo) SelectAll() ([]models.User, error) {
	query := fmt.Sprintf("SELECT %s FROM users", userColumns)
	rows, err := ur.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.Email, &u.Login, &u.Password)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

func (ur *UserRepo) FindByID(id string) (*models.User, error) {
	query := fmt.Sprintf("SELECT %s FROM users WHERE id=$1", userColumns)
	row := ur.storage.db.QueryRow(query, id)
	var user models.User

	err := row.Scan(&user.ID, &user.Email, &user.Login, &user.Password)
	if err != nil {
		return nil, appError.NotFoundError(err, "cannot find user")
	}
	return &user, nil
}

func (ur *UserRepo) FindByLogin(login string) (*models.User, error) {
	query := fmt.Sprintf("SELECT %s FROM users WHERE login=$1", userColumns)
	row := ur.storage.db.QueryRow(query, login)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Login, &user.Password)
	if err != nil {
		return nil, appError.NotFoundError(nil, "cannot find user with this login")
	}
	return &user, nil
}
