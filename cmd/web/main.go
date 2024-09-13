package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"snippetbox/internal/models"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

type application struct {
    debug          bool
    logger         *slog.Logger
    templateCache  map[string]*template.Template
    sessionManager *scs.SessionManager
    user           models.UserModelInterface
    snippet        models.SnippetModelInterface
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    dbDriver := flag.String("driver", "mysql", "Database driver name")
    dsn := flag.String("dsn", "zzh:zzhpwd@tcp(localhost:3306)/zsnippetbox?parseTime=true", "Data source name")
    debug := flag.Bool("debug", false, "Enable debug mode")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(*dbDriver, *dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    defer db.Close()

    templateCache, err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    sessionManager := scs.New()
    sessionManager.Store = mysqlstore.New(db)
    sessionManager.Lifetime = 12 * time.Hour
    sessionManager.Cookie.Secure = true // Setting this means the cookie will only be sent by a user's web browser when an HTTPS connection is used.

    app := &application{
        debug:          *debug,
        logger:         logger,
        templateCache:  templateCache,
        sessionManager: sessionManager,
        user:           &models.UserModel{DB: db},
        snippet:        &models.SnippetModel{DB: db},
    }

    tlsConfig := &tls.Config{
        CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
    }

    srv := &http.Server{
        Addr:         *addr,
        Handler:      app.routes(),
        ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
        TLSConfig:    tlsConfig,
        IdleTimeout:  time.Minute,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    idleConnsClosed := make(chan struct{})
    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt)
        <-sigint

        // We received an interrupt or kill signal, shut down the HTTP server.
        if err := srv.Shutdown(context.Background()); err != nil {
            // Error from closing listeners, or context timeout:
            logger.Error(fmt.Sprintf("HTTP server Shutdown: %v", err.Error()))
        }

        close(idleConnsClosed)
    }()

    logger.Info("starting HTTP server", "addr", *addr)

    err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
    if err != nil {
        if errors.Is(err, http.ErrServerClosed) {
            logger.Error("HTTP server closed")
        } else {
            // Error starting or closing listener:
            logger.Error(fmt.Sprintf("HTTP server ListenAndServe: %v", err.Error()))
        }
    }

    <-idleConnsClosed
}
