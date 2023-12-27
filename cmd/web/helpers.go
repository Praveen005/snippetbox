package main

import(
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