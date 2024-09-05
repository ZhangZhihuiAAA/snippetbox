package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    snippets, err := app.snippet.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    data := app.newTemplateData(r)
    data.Snippets = snippets

    app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    snippet, err := app.snippet.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w, r)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    data := app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w, r, http.StatusOK, "view.html", data)
}

// Update our snippetCreateForm struct to include struct tags which tell the decoder how to map
// HTML form values into the different struct fields. So, for example, here we're telling the
// decoder to store the value from the HTML form input with the name "title" in the Title field.
// The struct tag `form:"-"` tells the decoder to completely ignore a field during decoding.
type snippetCreateForm struct {
    Title               string `form:"title"`
    Content             string `form:"content"`
    Expires             int    `form:"expires"`
    validator.Validator `form:"-"`
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)

    // Initialize a new createSnippetForm instance and pass it to the template.
    // Notice how this is also a great opportunity to set any default or
    // 'initial' values for the form.
    data.Form = snippetCreateForm{
        Expires: 365,
    }

    app.render(w, r, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    var form snippetCreateForm

    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form.CheckField(validator.NotEmpty(form.Title), "title", "This field cannot be empty.")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long.")
    form.CheckField(validator.NotEmpty(form.Content), "content", "This field cannot be empty.")
    form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365.")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
        return
    }

    id, err := app.snippet.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Use the Put() method to add a string value and the corresponding key to the session data.
    app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
