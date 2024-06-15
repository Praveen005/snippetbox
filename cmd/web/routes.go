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
	




	// Leave the static files route unchanged. We don't need the session middleware to act on this route
	// Because all it does is serve static files.
	fileserver := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))



	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	// Unprotected application routes using the "dynamic" middleware chain. // pg. 300
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)	
	// dynamic := alice.New(app.sessionManager.LoadAndSave)	


	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))



	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	// Means you can simply type in https://localhost:4000/snippet/create and create a snippet, 
	// Only authenticted users should be able to do that.
	// Notice, it's dynamic.Append below, we are inserting requireAuthentication middleware in between.
	// Because the 'protected' middleware chain appends to the 'dynamic' chain
	// the noSurf middleware will also be used on the three routes below too.
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))


	// Using justinas/alice package to chain middleware
	// return alice.New(app.recoverPanic, app.logRequest, secureHeaders).Then(mux)
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}