package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) isAuthenticated(req *http.Request) bool {
	return app.session.Exists(req, "authenticatedUserID")
}

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

func (app *application) addDefaultData(data *templateData, req *http.Request) *templateData {
	if data == nil {
		data = &templateData{}
	}

	data.CSRFToken = nosurf.Token(req)
	data.CurrentYear = time.Now().Year()
	data.Flash = app.session.PopString(req, "flash")
	data.IsAuthenticated = app.isAuthenticated(req)
	return data
}

func (app *application) render(resWriter http.ResponseWriter, req *http.Request, name string, data *templateData) {
	templateSet, ok := app.templateCache[name]
	if !ok {
		app.serverError(resWriter, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buffer := new(bytes.Buffer)
	// catch runtime errors by trial rendering template as stream
	err := templateSet.Execute(buffer, app.addDefaultData(data, req))

	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	buffer.WriteTo(resWriter)
}
