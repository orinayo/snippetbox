package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"orinayooyelade.com/snippetbox/pkg/forms"
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
	app.render(resWriter, req, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(resWriter http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		app.clientError(resWriter, http.StatusBadRequest)
		return
	}

	form := forms.New(req.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(resWriter, req, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(resWriter, err)
		return
	}

	app.session.Put(req, "flash", "Snippet successfully created!")

	http.Redirect(resWriter, req, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(resWriter http.ResponseWriter, req *http.Request) {
	app.render(resWriter, req, "signup.page.tmpl", &templateData{
		Form: forms.New(nil)})
}

func (app *application) signupUser(resWriter http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		app.clientError(resWriter, http.StatusBadRequest)
		return
	}

	form := forms.New(req.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(resWriter, req, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(resWriter, req, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(resWriter, err)
		}
		return
	}

	app.session.Put(req, "flash", "Your signup was successful. Please log in.")
	http.Redirect(resWriter, req, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(resWriter http.ResponseWriter, req *http.Request) {
	app.render(resWriter, req, "login.page.tmpl", &templateData{Form: forms.New(nil)})
}
func (app *application) loginUser(resWriter http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		app.clientError(resWriter, http.StatusBadRequest)
		return
	}

	form := forms.New(req.PostForm)
	form.Required("email", "password")
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(resWriter, req, "login.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(resWriter, req, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(resWriter, err)
		}
		return
	}

	app.session.Put(req, "authenticatedUserID", id)
	http.Redirect(resWriter, req, "/snippet/create", http.StatusSeeOther)
}
func (app *application) logoutUser(resWriter http.ResponseWriter, req *http.Request) {
	app.session.Remove(req, "authenticatedUserID")
	app.session.Put(req, "flash", "You've been logged out successfully!")
	http.Redirect(resWriter, req, "/", http.StatusSeeOther)
}
