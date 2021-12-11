package api

import (
	"encoding/json"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/pkg/logger"
	"net/http"
)

// Message Вспомогательная структура для формирования сообщений
type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func initHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// @Summary     Register
// @Description  create account
// @Tags         auth
// @ID create-account
// @Accept       json
// @Produce      json
// @Param        input body models.User true "account info"
// @Success      201  {object}   Message
// @Router       /auth/register [post]
func (api *API) postUserRegister(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("Post User Register POST /api/user/register ")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.InfoLogger.Println("Invalid json received from client")
		return appError.NewAppError(err, "Provided json is invalid", http.StatusBadRequest)
	}

	userAdded, err := api.service.User().Create(&user)
	if err != nil {
		logger.InfoLogger.Println("Troubles while accessing database table (users) err:", err)
		return err
	}

	msg := Message{
		StatusCode: http.StatusCreated,
		Message:    fmt.Sprintf("User {login:%s} successfully registered!", userAdded.Login),
		IsError:    false,
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(msg)
}

func (api *API) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	users, err := api.service.User().SelectAll()
	if err != nil {
		logger.InfoLogger.Println(err)
		return appError.NewAppError(err, err.Error(), http.StatusInternalServerError)
	}
	logger.InfoLogger.Println("Ger All Users GET /users")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(users)
}
