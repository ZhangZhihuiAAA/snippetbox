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

// Validator contains structures that hold validation errors.
type Validator struct {
    NonFieldErrors []string  // Holds validation errors which are not related to a specific form field.
    FieldErrors    map[string]string  // Holds validation errors for form fields.
}

// Valid returns true if both NonFieldErrors and FieldErrors are empty.
func (v *Validator) Valid() bool {
    return len(v.NonFieldErrors) == 0 && len(v.FieldErrors) == 0
}

// AddNonFieldError adds error messages to the NonFieldErrors slice.
func (v *Validator) AddNonFieldError(message string) {
    v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// AddFieldError adds an error message to the FieldErrors map (so long as
// no entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
    // Note: We need to initialize the map first if it isn't already
    // initialized.
    if v.FieldErrors == nil {
        v.FieldErrors = make(map[string]string)
    }

    if _, exists := v.FieldErrors[key]; !exists {
        v.FieldErrors[key] = message
    }
}

// CheckField adds an error message to the FieldErrors map only if a
// validation check is not `ok`.
func (v *Validator) CheckField(ok bool, key string, message string) {
    if !ok {
        v.AddFieldError(key, message)
    }
}

// NotEmpty returns true if a value is not an empty string.
func NotEmpty(value string) bool {
    return strings.TrimSpace(value) != ""
}

// MinChars returns true if a value contains at least n characters.
func MinChars(value string, n int) bool {
    return utf8.RuneCountInString(value) >= n
}

// MaxChars returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
    return utf8.RuneCountInString(value) <= n
}

// PermittedValue returns true if a value is in a list of specific permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
    return slices.Contains(permittedValues, value)
}

// Match returns true if a value matches a provided compiled regular expression pattern.
func Match(value string, rx *regexp.Regexp) bool {
    return rx.MatchString(value)
}
