package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Praveen005/snippetbox/internal/models"
	"github.com/Praveen005/snippetbox/internal/validator"

	"github.com/julienschmidt/httprouter"
)



// Remove the explicit FieldErrors struct field and instead embed the Validator
// type. Embedding this means that our snippetCreateForm "inherits" all the
// fields and methods of our Validator type (including the FieldErrors field).
type snippetCreateForm struct{
	Title 			string
	Content 		string
	Expires 		int
	validator.Validator  // struct embedding: check the blog
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
	}

	// Because the Validator type is embedded(cheeck struct embedding blog) by the snippetCreateForm struct,
	// we can call CheckField() directly on it to execute our validation checks.
	// CheckField() will add the provided key and error message to the
	// FieldErrors map if the check does not evaluate to true. For example, in
	// the first line here we "check that the form.Title field is not blank". In
	// the second, we "check that the form.Title field has a maximum character
	// length of 100" and so on.
	// I was confused a bit about validator in validor.xyz() below, this validator is the package name. 
	// yes the package name, you can see nowhere I defined `validator`. We can use it this way,
	// beacuse we have embedded Validator struct from validator package, in our snippetCreateForm struct.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")


	// Use the Valid() method to see if any of the checks failed. If they did,
	// then re-render the template passing in the form in the same way as
	// before.

	if !form.Valid() {
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