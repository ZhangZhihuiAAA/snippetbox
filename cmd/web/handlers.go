package main

import (
	"fmt"
	"html/template"
	"net/http"
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

    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Save a new snippet..."))
}
