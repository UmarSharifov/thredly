package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/profile", app.profile)
	mux.HandleFunc("/moreuser", app.moreuser)
	mux.HandleFunc("/recs", app.recs)
	mux.HandleFunc("/admin", app.adminPage)
	mux.HandleFunc("/report", app.report)
	mux.HandleFunc("/subscribe", app.subscribe)
	mux.HandleFunc("/events", app.events)
	mux.HandleFunc("/eventscreate", app.createEvent)
	mux.HandleFunc("/createthred", app.createThred)
	mux.HandleFunc("/auth", app.auth)
	mux.HandleFunc("/logout", app.logout)
	mux.HandleFunc("/registration", app.registration)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/profile/detail", app.profileDetail)

	fileServer := http.FileServer(http.Dir("..\\..\\ui\\static\\"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
