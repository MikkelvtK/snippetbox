package main

import (
	"context"
	"net/http"
)

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := a.sessionManager.GetInt(r.Context(), "authenticatedUserId")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := a.users.Exists(id)
		if err != nil {
			a.serverError(w, r, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
