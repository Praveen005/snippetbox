package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Praveen005/snippetbox/internal/models"
)

// Add a Form field with the type "any".
// We’ll use this Form field to pass the validation errors and previously submitted data back to
// the template when we re-display the form.

type templateData struct{
	CurrentYear 	int
	Snippet     	*models.Snippet
	Snippets   		[]*models.Snippet
	Form			any
	Flash 			string // Add a Flash field to the templateData struct.
}


// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
    // Format the time in the desired format
    return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}


func newTemplateCache()(map[string]*template.Template, error){
	// Initializa a new map to act as a cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get the slice of all the filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". This will essentially give
	// us a slice of all the filepaths of our application 'page' templates
	// like: [ui/html/pages/home.tmpl  ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil{
		return nil, err
	}

	// Loop through the pages filepath one by one
	for _, page := range pages{
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil{
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil{
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		// Every page will have a base and partials,
		// So for every page(we are inside loop), we parse base, partials and the respective page
		// we have one home page, which displays the latest snippets created, other page is views page
		// which displays id specific snippet.
		// So, basically we are catching two pages, HomePage and ViewPage.
		ts, err = ts.ParseFiles(page)
		if err != nil{
			return nil, err
		}

		// Add the template set to the map as normal...
		cache[name] = ts
	}
	// return the map
	return cache, nil
}