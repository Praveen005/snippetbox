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
// ----------
// Update our snippetCreateForm struct to include struct tags which tell the
// decoder how to map HTML form values into the different struct fields. So, for
// example, here we're telling the decoder to store the value from the HTML form
// input with the name "title" in the Title field. The struct tag `form:"-"` 
// tells the decoder to completely ignore a field during decoding.

type snippetCreateForm struct {
	Title 					string 	`form:"title"`
	Content 				string 	`form:"content"`
	Expires 				int 	`form:"expires"`
	validator.Validator 			`form:"-"`
}

// Create a new userSignupForm struct.
type userSignupForm struct{
	Name 		string		`form:"name"`
	Email		string      `form:"email"`
	Password	string		`form:"password"`
	validator.Validator		`form:"-"`
}

// Create a new userLoginForm struct.
type userLoginForm struct{
	Email			string		`form:"email"`
	Password		string		`form:"password"`
	validator.Validator			`form:"-"`
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

	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.
	flash := app.sessionManager.PopString(r.Context(), "flash")


	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Pass the flash message to the template.
	data.Flash = flash
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
	
	// Declare a new empty instance of the snippetCreateForm struct.
	var form snippetCreateForm


	err := app.decodePostForm(r, &form)
	if err != nil{
		app.clientError(w, http.StatusBadRequest)
		return
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

	// Use the Put() method to add a string value ("Snippet successfully 
	// created!") and the corresponding key ("flash") to the session data.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")


	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}


func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func(app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare an zero-valued instance of our userSignupForm struct.
	var form userSignupForm


	// Parse the form data into the userSignupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRx), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be atleat 8 characters long")

	// If there are any errors, redisplay the signup form along with a 422
	// status code.
	if !form.Valid(){
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}


	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil{
		if errors.Is(err, models.ErrDuplicateEmail){
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else{
			app.serverError(w, err)
		}
		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {

	form := userLoginForm{}
	data := app.newTemplateData(r)
	data.Form = form
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func(app *application) userLoginPost(w http.ResponseWriter, r *http.Request){
	// Decode the form data into the userLoginForm struct.
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil{
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRx), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid(){
		data := app.newTemplateData(r)
		data.Form= form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}


	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil{
		if errors.Is(err, models.ErrInvalidCredentials){
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		}else{
			app.serverError(w, err)
		}
		return
	}

	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the 
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil{
		app.serverError(w, err)
		return
	}

	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// /Redirect the user to the create snippet page
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}