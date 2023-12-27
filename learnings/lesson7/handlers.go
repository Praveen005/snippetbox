package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)


func (app *application)home(w http.ResponseWriter, r* http.Request){
	if r.URL.Path != "/"{
		// http.NotFound(w, r)
		app.notFound(w) //from the helpers.go file
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil{
		// app.errorLog.Print(err.Error())
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverError(w, err) //serverError() function from helpers.go
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil{
		// app.errorLog.Print(err.Error())
		// http.Error(w, "Internal Server Error", 500)

		app.serverError(w, err) //serverError() function from helpers.go

		return
	}
} 

func (app *application) snippetView(w http.ResponseWriter, r * http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1{
		// http.NotFound(w, r) //sends a 404 response
		app.notFound(w) //notFound() function from helpers.go
		return
	}

	fmt.Fprintf(w, "Display the specific snippet with ID %d ...", id) 
}

func(app * application) snippetCreate(w http.ResponseWriter, r* http.Request){
	if r.Method != http.MethodPost{
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Method not Allowed", 405)
		// http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed) //from the helpers.go
		return
	}

	w.Write([]byte("Create a new snippet..."))
}