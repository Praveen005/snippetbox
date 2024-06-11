package main

import(
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
// http.ServeMux too implements the http.Handler interface, don't get confused.
func(app *application) routes() http.Handler{
	// Initialize the router
	router := httprouter.New()


	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Allowed responses by setting
	// router.MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	




	// Update the pattern for the route for the static files.
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))

	// And then create the routes using the appropriate methods, patterns and 
	// handlers.
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)  // Display(gets you) a HTML form for creating a new snippet
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)


	// Using justinas/alice package to chain middleware
	// return alice.New(app.recoverPanic, app.logRequest, secureHeaders).Then(mux)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}