package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func CreateUser(username string, email string, password string) map[string]string {
	db := turso.GetDB()
	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
		Score:    0,
	}

	db.Create(&user)

	data := SanitizeUserData(&user)

	return data
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

	res := SanitizeUserData(&user)

	return WriteJSON(w, http.StatusOK, res)
}

func Login(username string, password string) (map[string]string, error) {
	var user models.User
	var err error
	db := turso.GetDB()

	result := db.First(&user, "username = ?", username)

	if result.Error != nil {
		err = result.Error
	}

	data := SanitizeUserData(&user)

	return data, err
}

type LoginBody struct {
	Username string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	var login LoginBody
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, login)
}

func UpdateUserScore(id uuid.UUID, newScore int) (map[string]string, error) {
	var err error

	db := turso.GetDB()

	user, e := GetUser(id)

	if e != nil {
		err = e
	}

	var newUser *models.User

	newUser = user
	newUser.Score = newScore

	result := db.Save(&newUser)

	if result.Error != nil {
		err = result.Error
	}

	user, e = GetUser(id)
	updatedInfo := SanitizeUserData(user)
	return updatedInfo, err
}

func SanitizeUserData(user *models.User) map[string]string {
	data := make(map[string]string)

	data["username"] = user.Username
	data["email"] = user.Email
	data["score"] = strconv.Itoa(user.Score)

	return data
}

func UpdateUserScoreHandler(w http.ResponseWriter, r *http.Request) error {
	var UpdateScoreData models.UpdateUserStruct
	err := json.NewDecoder(r.Body).Decode(&UpdateScoreData)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: err.Error()})
	}

	_, e := GetUser(UpdateScoreData.ID)

	if e != nil {
		return WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: e.Error()})
	}

	res, err := UpdateUserScore(UpdateScoreData.ID, UpdateScoreData.Score)

	return WriteJSON(w, http.StatusOK, res)
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
	// cleanUser := make(map[string]string)

	id, e := uuid.Parse(idParam)
	if e != nil {
		return WriteJSON(w, http.StatusBadRequest, models.ApiError{Error: e.Error()})
	}
	user, err := GetUser(id)
	if err != nil {
		return err
	}

	cleanData := SanitizeUserData(user)
	return WriteJSON(w, http.StatusOK, cleanData)
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
