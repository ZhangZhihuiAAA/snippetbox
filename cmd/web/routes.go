package main

import (
	"net/http"
	"zsnippetbox/ui"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
    mux := http.NewServeMux()

    mux.Handle("GET /static/", http.FileServerFS(ui.Files))

    mux.HandleFunc("GET /ping", ping)

    // Unprotected routes using the "dynamic" middleware chain.
    dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF, app.authenticate)

    mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
    mux.Handle("GET /about", dynamic.ThenFunc(app.about))
    mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
    mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
    mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
    mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
    mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

    // Protected (authenticated-only) routes using the "protected" middleware chain which includes 
    // the requireAuthentication middleware.
    protected := dynamic.Append(app.requireAuthentication)

    mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
    mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
    mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
    mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPasswordUpdatePost))
    mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
    mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))

    standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

    return standard.Then(mux)
}