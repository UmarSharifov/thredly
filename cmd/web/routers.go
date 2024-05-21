package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/profile", app.profile)
	mux.HandleFunc("/createthred", app.createThred)
	mux.HandleFunc("/auth", app.auth)
	mux.HandleFunc("/logout", app.logout) // Добавьте этот маршрут
	mux.HandleFunc("/registration", app.registration)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/profile/detail", app.profileDetail)

	fileServer := http.FileServer(http.Dir("..\\..\\ui\\static\\"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
