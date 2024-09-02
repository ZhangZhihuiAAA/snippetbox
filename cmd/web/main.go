package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")

    // Importantly, we use the flag.Parse() function to parse the command line flag. 
    // This reads in the command-line flag value and assigns it to th addr variable. 
    // You need to call this *before* you use the addr variable, otherwise it will 
    // always contain the default value of ":4000". If any errors are encountered 
    // during parsing the application will be terminated.
    flag.Parse()

    mux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

    mux.HandleFunc("GET /{$}", home)
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    // The value returned from the flag.String() function is a pointer to the flag 
    // value, not the value itself.
    log.Printf("starting server on %s", *addr)

    err := http.ListenAndServe(*addr, mux)
    log.Fatal(err)
}