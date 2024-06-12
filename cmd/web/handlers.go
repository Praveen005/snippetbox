package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Praveen005/snippetbox/internal/models"

	"github.com/julienschmidt/httprouter"
)



// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
type snippetCreateForm struct{
	Title 			string
	Content 		string
	Expires 		int
	FieldErrors 	map[string]string
}


func (app *application)home(w http.ResponseWriter, r* http.Request){

	// Because httprouter matches the "/" path exactly, we can now remove the
	// manual check of r.URL.Path != "/" from this handler.

	// if r.URL.Path != "/"{
	// 	app.notFound(w)
	// 	return
	// }


	snippets, err := app.snippets.Latest()
	if err != nil{
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Pass the data to the render() helper as normal.
	app.render(w, http.StatusOK, "home.tmpl", data)
} 




func (app *application) snippetView(w http.ResponseWriter, r * http.Request){

	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. We'll talk about request context
	// in detail later in the book, but for now it's enough to know that you can
	// use the ParamsFromContext() function to retrieve a slice containing these
	// parameter names and values like so:
	params := httprouter.ParamsFromContext(r.Context())


	// id, err := strconv.Atoi(r.URL.Query().Get("id"))

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil || id < 1{
		// http.NotFound(w, r) //sends a 404 response
		app.notFound(w) //notFound() function from helpers.go
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(w)
		}else{
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	
	// Use the render helper
	app.render(w, http.StatusOK, "view.tmpl", data)

}

// Add a new snippetCreate handler, which for now returns a placeholder
// response. We'll update this shortly to show a HTML form.
func(app *application) snippetCreate(w http.ResponseWriter, r *http.Request){
	data := app.newTemplateData(r)


	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the 
	// snippet expiry to 365 days.
	// we do this because, till there is no field error, FieldErrors is a nil map, and hence
	// {{with .Form.FieldErrors.title}} throws error.
	// By this workaround, we will initialize te other fields with their zero values.    
	data.Form = snippetCreateForm{
		Expires: 365,
	}



	// newTemplateCache function is called in main.go, it parses all the pages, how?
	// Ever page will have the base same(chassy will remain same, the upper body might change, based on which 
	// page you are viewing. different handlers renders different pages). Also, base.tmpl invokes nav.tmpl,
	// so nav will remain same on all the pages. So, all the pages have already been parsed when we start our 
	// program and stored in 'templateCache' as template state(ts), and these are ready to be executed by
	// render function. below are rendering "create.tmpl" for this handler function, which displays
	// HTML form. when render is called, it will search for the template state of this page, "create.tmpl"
	// in templateCache and then execute that template.
	app.render(w, http.StatusOK, "create.tmpl", data)
}



// Renamed previous snippetCreate() to snippetCreatePost to write to the database
func(app * application) snippetCreatePost(w http.ResponseWriter, r* http.Request){
	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper to 
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
		FieldErrors: map[string]string{},
	}

	// Check that the title value is not blank and is not more than 100
	// characters long. If it fails either of those checks, add a message to the
	// errors map using the field name as the key.
	if strings.TrimSpace(form.Title) == ""{
		form.FieldErrors["title"] = "This field cannot be blank"
	}else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// check that content value is not blank
	if strings.TrimSpace(form.Content) == ""{
		form.FieldErrors["content"] = "This field cannot be blank"
	}


	// Check the expires value matches one of the permitted values (1, 7 or
	// 365).
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// If there are any validation errors re-display the create.tmpl template,
	// passing in the snippetCreateForm instance as dynamic data in the Form 
	// field. Note that we use the HTTP status code 422 Unprocessable Entity 
	// when sending the response to indicate that there was a validation error.
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}


	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil{
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}