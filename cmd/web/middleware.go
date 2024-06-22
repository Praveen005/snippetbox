package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}


func (app *application) recoverPanic(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		//  of a panic as Go unwinds the stack).

		defer func(){
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil{
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")

				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func(app *application) requireAuthentication(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		// Isn't it? why would you want an unAuthenticated request to go down the chain? You won't
		if !app.isAuthenticated(r){
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		// Its effect is specific to the response of the /snippet/create page.
		// It does not affect the caching behavior of other pages within your application.
		// its effect is tied to the specific request/response cycle 
		// initiated when the user accesses and interacts with the /snippet/create page.
		w.Header().Add("Cache-Control", "no-store")


		// And call the next handler in the chain
		next.ServeHTTP(w, r)

	})
}



// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly attributes set.
// Notice, the next handler is being wrapped inside noSurf middleware ensuring the application of CSRF 
// safeguards on all the requests that proceeds via this. 
// Remember we only need to apply on unsafe methods(non-GET/HEAD/OPTIONS/TRACE)
// Point to remember:  If per-session token implementations occur after the initial generation of a token, 
// the value is stored in the session(keep in mind this is why while chaining the middleware, we need the session manager middleware before the csrf one, your doubt stand cleared :p) and is used for each subsequent request until the session expires. : [ref](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#synchronizer-token-pattern): 
func noSurf(next http.Handler) http.Handler{
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})
	return csrfHandler
}


func(app *application)authenticate(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticatedUserID value from the session using the
		// GetInt() method. This will return the zero value for an int (0) if no
		// "authenticatedUserID" value is in the session -- in which case we
		// call the next handler in the chain as normal and return.
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0{
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise, we check to see if a user with that ID exists in our
		// database.
		exists, err := app.users.Exists(id)
		if err != nil{
			app.serverError(w, err)
			return
		}


		// If a matching user is found, we know that the request is
		// coming from an authenticated user who exists in our database. We
		// create a new copy of the request (with an isAuthenticatedContextKey
		// value of true in the request context) and assign it to r.
		if exists{
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}