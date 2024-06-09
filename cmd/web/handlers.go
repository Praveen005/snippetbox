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

	//It executes a parsed HTML template with provided data(nil here), generating the final output by merging the template's structure with dynamic content.
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

	//writing in the response body, will be displayed on the web page
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

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil{
		app.serverError(w, err)
		return
	}
	
	//Redirect the user to relevent page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}