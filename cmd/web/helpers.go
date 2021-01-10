package main

import (
	"bytes"
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

func (app *application) render(resWriter http.ResponseWriter, req *http.Request, name string, templateData *templateData) {
	templateSet, ok := app.templateCache[name]
	if !ok {
		app.serverError(resWriter, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buffer := new(bytes.Buffer)

	err := templateSet.Execute(buffer, templateData)

	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	buffer.WriteTo(resWriter)
}
