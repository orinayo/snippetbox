package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	resWriter.Write([]byte("Create a new snippet"))
}

func (app *application) createSnippet(resWriter http.ResponseWriter, req *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	http.Redirect(resWriter, req, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
