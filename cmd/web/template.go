package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"snippetbox/internal/models"
	"snippetbox/ui"
	"time"

	"github.com/justinas/nosurf"
)

type templateData struct {
    CurrentYear     int
    IsAuthenticated bool
    CSRFToken       string
    Flash           string
    Form            any
    Snippet         models.Snippet
    Snippets        []models.Snippet
    User            models.User
}

func humanDate(t time.Time) string {
    if t.IsZero() {
        return ""
    }

    // Convert the time to UTC before formatting it.
    return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
    "humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}

    // Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded filesystem which
    // match the pattern 'html/pages/*.html'. This essentially gives us a slice of all the 'page'
    // templates for the application.
    pages, err := fs.Glob(ui.Files, "html/pages/*.html")
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        name := filepath.Base(page)

        // Create a slice containing the filepath patterns for the template we want to parse.
        patterns := []string{
            "html/base.html",
            "html/partials/*.html",
            page,
        }

        // Use ParseFS() instead of ParseFiles() to parse the template files from the ui.Files
        // embedded filesystem.
        ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
        if err != nil {
            return nil, err
        }

        cache[name] = ts
    }

    return cache, nil
}

func (app *application) newTemplateData(r *http.Request) templateData {
    return templateData{
        CurrentYear:     time.Now().Year(),
        IsAuthenticated: app.isAuthenticated(r),
        CSRFToken:       nosurf.Token(r),
        Flash: app.sessionManager.PopString(r.Context(), "flash"),  // Add the flash message to the template data, if one exists.
    }
}
