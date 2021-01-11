package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"orinayooyelade.com/snippetbox/pkg/models"
)

func (app *application) home(resWriter http.ResponseWriter, req *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	app.render(resWriter, req, "home.page.tmpl", &templateData{Snippets: snippets})
}

func (app *application) showSnippet(resWriter http.ResponseWriter, req *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(resWriter)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(resWriter)
		} else {
			app.serverError(resWriter, err)
		}
		return
	}

	app.render(resWriter, req, "show.page.tmpl", &templateData{Snippet: snippet})
}

func (app *application) createSnippetForm(resWriter http.ResponseWriter, req *http.Request) {
	app.render(resWriter, req, "create.page.tmpl", nil)
}

func (app *application) createSnippet(resWriter http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		app.clientError(resWriter, http.StatusBadRequest)
		return
	}

	title := req.PostForm.Get("title")
	content := req.PostForm.Get("content")
	expires := req.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		app.render(resWriter, req, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData: req.PostForm,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	http.Redirect(resWriter, req, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
