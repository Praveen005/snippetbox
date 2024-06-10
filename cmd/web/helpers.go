package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

//The serverError helper writes an error and stack trace to the error log,
//Then sends a generic 500 Internal Server Error response to the user.

func(app * application) serverError(w http.ResponseWriter, err error){
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// app.errorLog.Print(trace)
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//The clientError helper sends a specific status code and corresponding description to the user. We'll use this later to send responses like 400
//"Bad Request" when there's  a problem with request that the user sent.

func (app *application) clientError(w http.ResponseWriter, status int){
	http.Error(w, http.StatusText(status), status)
}

//For consistency, we'll implement a notFound Helper. This is simply a convenience wrapper around clientError which sends a 404 Not Found response
// to the user

func(app *application) notFound(w http.ResponseWriter){
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData){
	// Retrieve the appropriate template set from the cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper
	// method that we made earlier and return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Inalize a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil{
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)


	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}