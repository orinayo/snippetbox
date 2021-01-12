package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Form struct which embeds a url.Values object (form Data)
// and an Errors field holding validation errors
type Form struct {
	url.Values
	Errors errors
}

// EmailRX uses the regexp.MustCompile() function to parse a pattern and compile a
// regular expression for sanity checking the format of an email address.
// This returns a *regexp.Regexp object, or panics in the event of an error.
// Doing this once at runtime, and storing the compiled regular expression
// object in a variable, is more performant than re-compiling the pattern with
// every request
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// New will initialize a custom Form struct.
// It takes form data as its parameter
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required will check that specific fields in the form
// are present and not blank
func (formData *Form) Required(fields ...string) {
	for _, field := range fields {
		value := formData.Get(field)
		if strings.TrimSpace(value) == "" {
			formData.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MinLength will check that a specific field in the form
// contains a minimum number of characters
func (formData *Form) MinLength(field string, min int) {
	value := formData.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < min {
		formData.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", min))
	}
}

// MaxLength will check that a specific field in the form
// contains a maximum number of characters
func (formData *Form) MaxLength(field string, max int) {
	value := formData.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > max {
		formData.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", max))
	}
}

// PermittedValues will check that a specific field in the form
// matches one of a set of specific permitted values
func (formData *Form) PermittedValues(field string, opts ...string) {
	value := formData.Get(field)
	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	formData.Errors.Add(field, "This field is invalid")
}

// MatchesPattern will check that a specific field in the form
// matches a regular expression
func (formData *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := formData.Get(field)
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		formData.Errors.Add(field, "This field is invalid")
	}
}

// Valid will return true if there are no errors
func (formData *Form) Valid() bool {
	return len(formData.Errors) == 0
}
