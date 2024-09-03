package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"snippetbox/internal/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    // Initialize a slice containing the paths to the two files. It's important 
    // to note that the file containing our base template must be the *first* 
    // file in the slice.
    files := []string{
        "./ui/html/base.html",
        "./ui/html/partials/nav.html",
        "./ui/html/pages/home.html",
    }

    // Use the template.ParseFiles() function to read the template files into a 
    // template set. If there's an error, we log the detailed error message, use the 
    // http.Error() function to send an Internal Server Error response to the user, 
    // and then return from the handler so no subsequent code is executed.
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, r, err)  // Use the serverError() helper.
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Then we use the Execute() method on the template set to write the template 
    // content as the response body. The last parameter to Execute() represents 
    // any dynamic data that we want to pass in, which for now we'll leave as nil.
    // err = ts.Execute(w, nil)

    // Use the ExecuteTemplate() method to write the content of the "base" 
    // template as the response body.
    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        app.serverError(w, r, err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
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

    // Initialize a slice containing the paths to the view.html file, 
    // plus the base layout and navigation partial that we made earlier.
    files := []string{
        "./ui/html/base.html",
        "./ui/html/partials/nav.html",
        "./ui/html/pages/view.html",
    }

    // Parse the template files.
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // And then execute them. Notice how we are passing in the snippet 
    // data (a models.Snippet struct) as the final parameter.
    err = ts.ExecuteTemplate(w, "base", snippet)
    if err != nil {
        app.serverError(w, r, err)
    }
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    title := "0 snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
    expires := 7

    id, err := app.snippet.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Redirect the user to the relevant page for the snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
