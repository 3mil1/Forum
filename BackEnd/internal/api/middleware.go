package api

import (
	"context"
	"errors"
	"fmt"
	"forum/internal/appError"
	"forum/internal/models"
	"forum/pkg/logger"
	"net/http"
	"strings"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var appErr *appError.AppError

		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				logger.InfoLogger.Println(appErr.Message, ":", appErr.Err)
				w.WriteHeader(appErr.StatusCode)
				w.Write(appErr.Marshal())
				return
			}
			logger.ErrorLogger.Println("Unhandled error occurred: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(appError.SystemError(err).Marshal())
		}
	}
}

type userContext struct {
	userID string
	login  string
	email  string
}

func (api *API) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		fmt.Println(c)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		xs := strings.SplitN(c.Value, "|", 2)
		if len(xs) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		cCode := xs[0]
		uLogin := xs[1]
		var user *models.User
		if user, err = api.service.Session().Check(cCode, uLogin); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// set context
		ctx := context.WithValue(r.Context(), "values", userContext{userID: user.ID, login: user.Login, email: user.Email})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CorsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, *")
		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
		}
	})
}
