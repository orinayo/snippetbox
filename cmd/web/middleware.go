package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		resWriter.Header().Set("X-XSS-Protection", "1; mode=block")
		resWriter.Header().Set("X-Frame-Options", "deny")

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
