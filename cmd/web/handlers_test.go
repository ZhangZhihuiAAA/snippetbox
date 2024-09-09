package main

import (
	"net/http"
	"net/url"
	"snippetbox/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
    app := newTestApplication(t)

    ts := newTestServer(t, app.routes())
    defer ts.Close()

    code, _, body := ts.get(t, "/ping")

    assert.Equal(t, code, http.StatusOK)
    assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
    // Create a new instance of our application struct which uses the mocked dependencies.
    app := newTestApplication(t)

    // Establish a new test server for running end-to-end tests.
    ts := newTestServer(t, app.routes())
    defer ts.Close()

    tests := []struct {
        name       string
        urlPath    string
        expectCode int
        expectBody string
    }{
        {
            name:       "Valid ID",
            urlPath:    "/snippet/view/1",
            expectCode: http.StatusOK,
            expectBody: "An old silent pond...",
        },
        {
            name:       "Non-existent ID",
            urlPath:    "/snippet/view/2",
            expectCode: http.StatusNotFound,
        },
        {
            name:       "Negative ID",
            urlPath:    "/snippet/view/-1",
            expectCode: http.StatusNotFound,
        },
        {
            name:       "Decimal ID",
            urlPath:    "/snippet/view/1.23",
            expectCode: http.StatusNotFound,
        },
        {
            name:       "String ID",
            urlPath:    "/snippet/view/foo",
            expectCode: http.StatusNotFound,
        },
        {
            name:       "Empty ID",
            urlPath:    "/snippet/view/",
            expectCode: http.StatusNotFound,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            code, _, body := ts.get(t, tc.urlPath)

            assert.Equal(t, code, tc.expectCode)

            if tc.expectBody != "" {
                assert.StringContains(t, body, tc.expectBody)
            }
        })
    }
}

func TestUserSignup(t *testing.T) {
    app := newTestApplication(t)
    ts := newTestServer(t, app.routes())
    defer ts.Close()

    _, _, body := ts.get(t, "/user/signup")
    validCSRFToken := extractCSRFToken(t, body)

    const (
        validName     = "Bob"
        validPassword = "validPa$$word"
        validEmail    = "bob@example.com"
        formTag       = `<form action="/user/signup" method="POST" novalidate>`
    )

    tests := []struct {
        name          string
        userName      string
        userEmail     string
        userPassword  string
        csrfToken     string
        expectCode    int
        expectFormTag string
    }{
        {
            name: "Valid submission",
            userName: validName,
            userEmail: validEmail,
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            expectCode: http.StatusSeeOther,
        },
        {
            name: "Invalid CSRF Token",
            userName: validName,
            userEmail: validEmail,
            userPassword: validPassword,
            csrfToken: "wrongToken",
            expectCode: http.StatusBadRequest,
        },
        {
            name: "Empty name",
            userName: "",
            userEmail: validEmail,
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            expectCode: http.StatusUnprocessableEntity,
            expectFormTag: formTag,
        },
        {
            name: "Empty email",
            userName: validName,
            userEmail: "",
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            expectCode: http.StatusUnprocessableEntity,
            expectFormTag: formTag,
        },
        {
            name: "Empty password",
            userName: validName,
            userEmail: validEmail,
            userPassword: "",
            csrfToken: validCSRFToken,
            expectCode: http.StatusUnprocessableEntity,
            expectFormTag: formTag,
        },
        {
            name: "Invalid email",
            userName: validName,
            userEmail: "bob@example.",
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            expectCode: http.StatusUnprocessableEntity,
            expectFormTag: formTag,
        },
        {
            name: "Short password",
            userName: validName,
            userEmail: validEmail,
            userPassword: "pa$$",
            csrfToken: validCSRFToken,
            expectCode: http.StatusUnprocessableEntity,
            expectFormTag: formTag,
        },
        {
            name: "Duplicate email",
            userName: validName,
            userEmail: "dupe@example.com",
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            expectCode: http.StatusUnprocessableEntity,
            expectFormTag: formTag,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            form := url.Values{}
            form.Add("name", tc.userName)
            form.Add("email", tc.userEmail)
            form.Add("password", tc.userPassword)
            form.Add("csrf_token", tc.csrfToken)

            code, _, body := ts.postForm(t, "/user/signup", form)

            assert.Equal(t, code, tc.expectCode)

            if tc.expectFormTag != "" {
                assert.StringContains(t, body, tc.expectFormTag)
            }
        })
    }
}
