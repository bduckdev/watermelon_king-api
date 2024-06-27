package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"watermelon_king-api/models"

	"github.com/go-chi/chi/v5"
	jwt "github.com/golang-jwt/jwt"
	"github.com/markbates/goth/gothic"
)

func CreateJWT(user *models.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 25000,
		"userId":    user.ID,
	}

	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func WithAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling jwt auth middleware")

		tokenString := r.Header.Get("x-jwt-token")

		_, err := ValidateJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, models.ApiError{Error: "Invalid Token"})
			return
		}

		handlerFunc(w, r)
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected token signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func GetAuthCallbackFunction(w http.ResponseWriter, r *http.Request) error {
	var err error

	provider := chi.URLParam(r, "provider")
	fmt.Println("Provider:", provider)

	provider = "google"

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, e := gothic.CompleteUserAuth(w, r)

	if e != nil {
		err = e
	}

	fmt.Println(user)
	http.Redirect(w, r, "http://localhost:5173/", http.StatusFound)
	return err
}
