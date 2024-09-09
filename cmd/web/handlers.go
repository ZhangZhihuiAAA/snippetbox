package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
)

func ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}

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

func (app *application) about(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    app.render(w, r, http.StatusOK, "about.html", data)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
    userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

    user, err := app.user.Get(userID)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.Redirect(w, r, "/user/login", http.StatusSeeOther)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    data := app.newTemplateData(r)
    data.User = user

    app.render(w, r, http.StatusOK, "account.html", data)
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

type snippetCreateForm struct {
    Title               string     `form:"title"`
    Content             string     `form:"content"`
    Expires             int        `form:"expires"`
    validator.Validator `form:"-"` // The struct tag `form:"-"` tells the decoder to completely ignore a field during decoding.
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)

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

    app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type userSignupForm struct {
    Name                string `form:"name"`
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = userSignupForm{}
    app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
    var form userSignupForm

    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form.CheckField(validator.NotEmpty(form.Name), "name", "This field cannot be empty.")
    form.CheckField(validator.NotEmpty(form.Email), "email", "This field cannot be empty.")
    form.CheckField(validator.Match(form.Email, validator.EmailRX), "email", "This field must be a valid email address.")
    form.CheckField(validator.NotEmpty(form.Password), "password", "This field cannot be empty.")
    form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long.")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
        return
    }

    // Try to create a new user record in the database. If the email already exists then add an
    // error message to the form and re-display it.
    err = app.user.Insert(form.Name, form.Email, form.Password)
    if err != nil {
        if errors.Is(err, models.ErrDuplicateEmail) {
            form.AddFieldError("email", "Email address is already in use.")

            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
        } else {
            app.serverError(w, r, err)
        }

        return
    }

    // Otherwise add a confirmation flash message to the session confirming that their signup worked.
    app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

    // And redirect the user to the login page.
    http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = userLoginForm{}
    app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
    var form userLoginForm

    err := app.decodePostForm(r, &form)
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form.CheckField(validator.NotEmpty(form.Email), "email", "This field cannot be empty.")
    form.CheckField(validator.Match(form.Email, validator.EmailRX), "email", "This field must be a valid email address.")
    form.CheckField(validator.NotEmpty(form.Password), "password", "This field cannot be blank.")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
        return
    }

    id, err := app.user.Authenticate(form.Email, form.Password)
    if err != nil {
        if errors.Is(err, models.ErrInvalidCredentials) {
            form.AddNonFieldError("Email or password is incorrect!")

            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    // Use the RenewToken() method on the current session to change the session ID. It's a good 
    // practice to generate a new session ID when the authentication state or privilage level 
    // changes for the user (e.g. login and logout operations).
    err = app.sessionManager.RenewToken(r.Context())
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Add the ID of the current user to the session, so that they are now 'logged in'.
    app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

    // Use the PopString method to retrieve and remove a value from the session data in one step. 
    // If no matching key exists this will return the empty string.
    path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
    if path != "" {
        http.Redirect(w, r, path, http.StatusSeeOther)
        return
    }

    http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
    // Use the RenewToken() method on the current session to change the session ID again.
    err := app.sessionManager.RenewToken(r.Context())
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Remove the authenticatedUserID from the session data so that the user is 'logged out'.
    app.sessionManager.Remove(r.Context(), "authenticatedUserID")

    // Add a flash message to the session to confirm to the user that they've been logged out.
    app.sessionManager.Put(r.Context(), "flash", "You've logged out successfully!")

    // Redirect the user to the application home page.
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
