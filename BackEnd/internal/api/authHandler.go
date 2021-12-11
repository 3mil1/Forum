package api

import (
	"encoding/json"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/pkg/logger"
	"net/http"
	"time"
)

// @Summary      Login
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Router       /auth/login [post]
func (api *API) postToAuth(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("Post to Auth POST /api/auth/login")
	var userFromJson models.User
	err := json.NewDecoder(r.Body).Decode(&userFromJson)

	if err != nil {
		fmt.Println(r.Body)
		logger.InfoLogger.Println("Invalid json received from client", userFromJson)
		return appError.NewAppError(err, "Provided json is invalid", http.StatusBadRequest)
	}

	code, err := api.service.Session().NewSession(userFromJson.Login, userFromJson.Password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	c := http.Cookie{
		Name:    "session",
		Value:   code + "|" + userFromJson.Login,
		Expires: time.Now().AddDate(1, 0, 0),
		Path:    "/",
		//HttpOnly: true,
		//SameSite: http.SameSiteNoneMode,
	}

	logger.InfoLogger.Println("Setting cookie", c)
	http.SetCookie(w, &c)

	return nil
}

func (api *API) authMe(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("GET to Auth me /api/auth/me")

	val, _ := r.Context().Value("values").(userContext)

	//setting values from context
	var user models.User
	user.ID = val.userID
	user.Email = val.email
	user.Login = val.login

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}

func (api *API) logOut(w http.ResponseWriter, r *http.Request) error {
	initHeaders(w)
	logger.InfoLogger.Println("Delete to log out Delete /api/auth/delete")

	val, _ := r.Context().Value("values").(userContext)
	var user models.User
	user.ID = val.userID

	err := api.service.Session().End(user.ID)
	if err != nil {
		return err
	}
	return nil
}
