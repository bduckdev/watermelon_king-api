package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"watermelon_king-api/models"
	"watermelon_king-api/pkg/turso"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func CreateUser(username string, email string, password string) *models.User {
	db := turso.GetDB()

	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
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
	err = models.ValidateNewUser(&user)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}

	CreateUser(user.Username, user.Email, user.Password)

	return WriteJSON(w, http.StatusOK, user)
}

func UpdateUser(id uuid.UUID) (*models.User, error) {
	db := turso.GetDB()
	var err error

	user, e := GetUser(id)

	if e != nil {
		err = e
	}

	result := db.Save(&user)

	if result.Error != nil {
		err = result.Error
	}

	return user, err
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) error {
	var user models.User

	UpdateUser(id)

	return WriteJSON(w, http.StatusOK, user)
}

func GetUser(id uuid.UUID) (*models.User, error) {
	var user models.User
	var err error

	db := turso.GetDB()

	result := db.First(&user, "id = ?", id)

	if result.Error != nil {
		err = result.Error
	}

	return &user, err
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) error {
	idParam := chi.URLParam(r, "id")

	id, e := uuid.Parse(idParam)
	if e != nil {
		return WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: e.Error()})
	}
	user, err := GetUser(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func GetAllUserScores() (map[string]int, error) {
	var users []*models.User
	var err error
	cleanUsers := make(map[string]int)

	db := turso.GetDB()

	result := db.Find(&users)
	if result.Error != nil {
		err = result.Error
	}
	for _, el := range users {
		cleanUsers[el.Username] = el.Score
	}

	return cleanUsers, err
}

func GetAllUserScoresHandler(w http.ResponseWriter, r *http.Request) error {
	users, err := GetAllUserScores()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, users)
}
