package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
    var (
        method = r.Method
        uri    = r.URL.RequestURI()
        trace  = string(debug.Stack())
    )

    app.logger.Error(err.Error(), "method", method, "uri", uri)

    if app.debug {
        body := fmt.Sprintf("%s\n\n%s", err.Error(), trace)
        http.Error(w, body, http.StatusInternalServerError)
        return
    }

    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
    ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("the template %s does not exist", page)
        app.serverError(w, r, err)
        return
    }

    buf := new(bytes.Buffer)

    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    w.WriteHeader(status)
    buf.WriteTo(w)
}

func (app *application) isAuthenticated(r *http.Request) bool {
    isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
    if !ok {
        return false
    }

    return isAuthenticated
}

func (app *application) decodePostForm(r *http.Request, varForm any) error {
    err := r.ParseForm()
    if err != nil {
        return err
    }

    err = form.NewDecoder().Decode(varForm, r.PostForm)
    if err != nil {
        // If we try to use an invalid target destination, the Decode() method will return an
        // error with the type *form.InvalidDecoderError. We use errors.As() to check for this
        // and raise a panic rather than returning the error.
        var invalidDecoderError *form.InvalidDecoderError

        if errors.As(err, &invalidDecoderError) {
            panic(err)
        }

        return err
    }

    return nil
}
