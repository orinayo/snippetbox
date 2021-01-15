package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	resRecorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our secureHeaders
	// middleware, which writes a 200 status code and "OK" response body.
	next := http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		resWriter.Write([]byte("OK"))
	})

	// Pass the mock HTTP handler to our secureHeaders middleware. Because
	// secureHeaders *returns* a http.Handler we can call its ServeHTTP()
	// method, passing in the http.ResponseRecorder and dummy http.Request to
	// execute it.
	secureHeaders(next).ServeHTTP(resRecorder, req)

	// Call the Result() method on the http.ResponseRecorder to get the results
	// of the test.
	res := resRecorder.Result()

	// Check that the middleware has correctly set the X-Frame-Options header
	// on the response.
	frameOptions := res.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	// Check that the middleware has correctly set the X-XSS-Protection header
	// on the response.
	xssProtection := res.Header.Get("X-XSS-Protection")
	if res.StatusCode != http.StatusOK {
		t.Errorf("want %q; got %q", "mode-block", xssProtection)
	}

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, res.StatusCode)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
