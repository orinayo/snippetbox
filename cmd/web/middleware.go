package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"orinayooyelade.com/snippetbox/pkg/models"
)

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		resWriter.Header().Set("X-XSS-Protection", "1; mode=block")
		resWriter.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(resWriter, req)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		exists := app.session.Exists(req, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(resWriter, req)
			return
		}

		user, err := app.users.Get(app.session.GetInt(req, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(req, "authenticatedUserID")
			next.ServeHTTP(resWriter, req)
			return
		} else if err != nil {
			app.serverError(resWriter, err)
			return
		}

		ctx := context.WithValue(req.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(resWriter, req.WithContext(ctx))
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		if !app.isAuthenticated(req) {
			app.session.Put(req, "redirectPathAfterLogin", req.URL.Path)
			http.Redirect(resWriter, req, "/user/login", http.StatusSeeOther)
			return
		}

		resWriter.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(resWriter, req)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", req.RemoteAddr, req.Proto, req.Method, req.URL.RequestURI())
		next.ServeHTTP(resWriter, req)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				resWriter.Header().Set("Connection", "close")
				app.serverError(resWriter, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(resWriter, req)
	})
}
