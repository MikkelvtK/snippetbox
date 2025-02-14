package main

import (
	"net/http"

	"github.com/MikkelvtK/snippetbox/ui"
	"github.com/justinas/alice"
)

func (a *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(a.sessionManager.LoadAndSave, noSurf, a.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(a.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(a.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(a.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(a.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(a.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(a.userLoginPost))

	protected := dynamic.Append(a.requireAuthentication)

	mux.Handle("POST /user/logout", protected.ThenFunc(a.userLogoutPost))
	mux.Handle("GET /snippet/create", protected.ThenFunc(a.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(a.snippetCreatePost))

	standard := alice.New(a.recoverPanic, a.logRequest, commonHeaders)

	return standard.Then(mux)
}
