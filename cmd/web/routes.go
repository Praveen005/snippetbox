package main

import(
	"net/http"
)

//The routes() method returns a servemux containing our application routes
func(app *application) routes() *http.ServeMux{
	mux := http.NewServeMux()

	//http.FileServer creates a simple file server that serves static files from a specified directory. 
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}