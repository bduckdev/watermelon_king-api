package handlers

import (
	"encoding/json"
	"net/http"
	"watermelon_king-api/models"
	"watermelon_king-api/pkg/turso"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func NewHandlerFunc(f models.ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: err.Error()})
		}
	}
}

func CreateUser(username string, email string) *models.User {
	db := turso.GetDB()

	user := models.User{
		Username: username,
		Email:    email,
		Score:    0,
	}

	db.Create(&user)

	return &user
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) error {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: err.Error()})
	}
	CreateUser(user.Username, user.Email)

	return WriteJSON(w, http.StatusOK, user)
}

func GetUser() *models.User {
	return &models.User{
		Username: "deez nuts",
		Email:    "deez nuts",
		Score:    69,
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) error {
	user := GetUser()

	return WriteJSON(w, http.StatusOK, user)
}
