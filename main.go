package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"watermelon_king-api/handlers"
	"watermelon_king-api/pkg/auth"
	"watermelon_king-api/pkg/turso"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func main() {
	Init()
	auth.NewAuth()
	r := chi.NewRouter()
	port := os.Getenv("PORT")
	r.Use(middleware.Logger)
	r.Post("/users", handlers.NewHandlerFunc(handlers.CreateUserHandler))
	r.Get("/users/{id}", handlers.NewHandlerFunc(handlers.GetUserHandler))
	r.Post("/users/{}", handlers.NewHandlerFunc(handlers.UpdateUserScoreHandler))
	r.Get("/users", handlers.NewHandlerFunc(handlers.GetAllUserScoresHandler))
	r.Get("/auth/{provider}/callback", handlers.NewHandlerFunc(handlers.GetAuthCallbackFunction))
	r.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	r.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
			t, _ := template.New("foo").Parse(userTemplate)
			t.Execute(res, gothUser)
		} else {
			gothic.BeginAuthHandler(res, req)
		}
	})

	fmt.Println("Server listening on port", port)
	http.ListenAndServe(port, r)
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	turso.Init()
}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
