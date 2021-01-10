package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(resWriter http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		app.notFound(resWriter)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	templateSet, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	err = templateSet.Execute(resWriter, nil)
	if err != nil {
		app.serverError(resWriter, err)
	}
}

func (app *application) showSnippet(resWriter http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(resWriter)
		return
	}
	fmt.Fprintf(resWriter, "Display a specific snippet with ID %d", id)
}

func (app *application) createSnippet(resWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		resWriter.Header().Set("Allow", http.MethodPost)
		app.clientError(resWriter, http.StatusMethodNotAllowed)
		return
	}
	resWriter.Write([]byte("Create a new snippet"))
}
