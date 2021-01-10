package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(resWriter http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(resWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(resWriter http.ResponseWriter, status int) {
	http.Error(resWriter, http.StatusText(status), status)
}

func (app *application) notFound(resWriter http.ResponseWriter) {
	app.clientError(resWriter, http.StatusNotFound)
}
