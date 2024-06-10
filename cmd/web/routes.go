package main

import(
	"net/http"
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


	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	return secureHeaders(mux)
}