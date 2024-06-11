package main

import(
	"net/http"

	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
// http.ServeMux too implements the http.Handler interface, don't get confused.
func(app *application) routes() http.Handler{
	mux := http.NewServeMux()

	//http.FileServer creates a simple file server that serves static files from a specified directory. 
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)


	// Wrap the existing chain with the logRequest middleware.
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	

	// Using justinas/alice package to chain middleware
	// return alice.New(app.recoverPanic, app.logRequest, secureHeaders).Then(mux)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}