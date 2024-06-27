package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "RANDOMAF"
	MaxAge = 86400 * 30
	IsProd = false
)

func NewAuth() {
	googleClientId := os.Getenv("GOOGLE_KEY")

	googleSecret := os.Getenv("GOOGLE_SECRET")

	store := sessions.NewCookieStore([]byte(key))

	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleSecret, "http://localhost:3001/auth/google/callback"),
	)
}
