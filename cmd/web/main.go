package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// application struct holds the application-wide dependencies for the web application.
type application struct {
    logger        *slog.Logger
    snippet       *models.SnippetModel
    templateCache map[string]*template.Template
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    dbDriver := flag.String("dbdriver", "mysql", "Database driver name")
    dsn := flag.String("dsn",
        "zeb:zebpwd@tcp(localhost:3306)/snippetbox?parseTime=true",
        "MySQL data source name")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(*dbDriver, *dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    defer db.Close()

    // Initialize a new template cache.
    templateCache, err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    app := &application{
        logger:        logger,
        snippet:       &models.SnippetModel{DB: db},
        templateCache: templateCache,
    }

    logger.Info("starting server", "addr", *addr)

    err = http.ListenAndServe(*addr, app.routes())
    logger.Error(err.Error())
    os.Exit(1)
}

func openDB(driverName string, dsn string) (*sql.DB, error) {
    db, err := sql.Open(driverName, dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}
