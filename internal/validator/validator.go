package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// Use the regexp.MustCompile() function to parse a regular expression pattern for sanity checking
// the format of an email address. This returns a pointer to a 'compiled' regexp.Regexp type, or
// panics in the event of an error. Parsing this pattern once at startup and storing the compiled
// *regexp.Regexp in a variable is more performant than re-parsing the pattern each time we need
// it.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// NotEmpty reports whether s is empty after being trimmed.
func NotEmpty(s string) bool {
    return strings.TrimSpace(s) != ""
}

// MinChars reports whether s contains at least n characters.
func MinChars(s string, n int) bool {
    return utf8.RuneCountInString(s) >= n
}

// MaxChars reports whether s contains no more than n characters.
func MaxChars(s string, n int) bool {
    return utf8.RuneCountInString(s) <= n
}

// PermittedValue reports whether v is one of permittedValues.
func PermittedValue[T comparable](v T, permittedValues ...T) bool {
    return slices.Contains(permittedValues, v)
}

// Match reports whether the string (s) and the compiled regular expression pattern (rx) match.
func Match(s string, rx *regexp.Regexp) bool {
    return rx.MatchString(s)
}

// Validator contains structures that hold validation errors.
type Validator struct {
    NonFieldErrors []string  // Holds validation errors which are not related to a specific form field.
    FieldErrors    map[string]string  // Holds validation errors for form fields.
}

// Valid reports whether both NonFieldErrors and FieldErrors of v are empty.
func (v *Validator) Valid() bool {
    return len(v.NonFieldErrors) == 0 && len(v.FieldErrors) == 0
}

// AddNonFieldError adds an error message to the NonFieldErrors slice of v.
func (v *Validator) AddNonFieldError(msg string) {
    v.NonFieldErrors = append(v.NonFieldErrors, msg)
}

// AddFieldError adds an error message to the FieldErrors map of v so long as no entry already exists 
// for the given field.
func (v *Validator) AddFieldError(field, msg string) {
    if v.FieldErrors == nil {
        v.FieldErrors = make(map[string]string)
    }

    if _, ok := v.FieldErrors[field]; !ok {
        v.FieldErrors[field] = msg
    }
}

// CheckField adds an error message to the FieldErrors map of v only if validation check of `field` is not `ok`.
func (v *Validator) CheckField(ok bool, field, msg string) {
    if !ok {
        v.AddFieldError(field, msg)
    }
}