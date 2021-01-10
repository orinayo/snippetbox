package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"orinayooyelade.com/snippetbox/pkg/models"
)

func (app *application) home(resWriter http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		app.notFound(resWriter)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	app.render(resWriter, req, "home.page.tmpl", &templateData{Snippets: snippets})
}

func (app *application) showSnippet(resWriter http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
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

func (app *application) createSnippet(resWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		resWriter.Header().Set("Allow", http.MethodPost)
		app.clientError(resWriter, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	http.Redirect(resWriter, req, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
