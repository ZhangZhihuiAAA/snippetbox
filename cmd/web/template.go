package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippetbox/internal/models"
	"snippetbox/ui"
	"time"
)

type templateData struct {
    CurrentYear     int
    Snippet         models.Snippet
    Snippets        []models.Snippet
    Form            any
    Flash           string
    IsAuthenticated bool
    CSRFToken       string
}

func humanDate(t time.Time) string {
    return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
    "humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}

    // Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded filesystem which 
    // match the pattern 'html/pages/*.html'. This essentially gives us a slice of all the 'page' 
    // templates for the application, just like before.
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
