package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	code, _, body := testServer.get(t, "/ping")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)

	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tableTest := range tests {
		t.Run(tableTest.name, func(t *testing.T) {
			code, _, body := testServer.get(t, tableTest.urlPath)
			if code != tableTest.wantCode {
				t.Errorf("want %d; got %d", tableTest.wantCode, code)
			}
			if !bytes.Contains(body, tableTest.wantBody) {
				t.Errorf("want body to contain %q", tableTest.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)

	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	_, _, body := testServer.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters)")},
		{"Duplicate email", "Bob", "andy@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := testServer.postForm(t, "/user/signup", form)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestCreateSnippetForm(t *testing.T) {
	app := newTestApplication(t)
	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := testServer.get(t, "/snippet/create")
		if code != http.StatusSeeOther {
			t.Errorf("want %d; got %d", http.StatusSeeOther, code)
		}
		if headers.Get("Location") != "/user/login" {
			t.Errorf("want %s; got %s", "/user/login", headers.Get("Location"))
		}
	})

	t.Run("Authenticated", func(t *testing.T) {
		// Authenticate the user...
		_, _, body := testServer.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "")
		form.Add("csrf_token", csrfToken)
		testServer.postForm(t, "/user/login", form)
		// Then check that the authenticated user is shown the create snippet form.
		code, _, body := testServer.get(t, "/snippet/create")
		if code != 200 {
			t.Errorf("want %d; got %d", 200, code)
		}
		formTag := "<form action='/snippet/create' method='POST'>"
		if !bytes.Contains(body, []byte(formTag)) {
			t.Errorf("want body %s to contain %q", body, formTag)
		}
	})
}
