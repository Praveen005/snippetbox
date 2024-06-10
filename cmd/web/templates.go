package main

import (
	"html/template"
	"path/filepath"

	"github.com/Praveen005/snippetbox/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct{
	Snippet    *models.Snippet
	Snippets   []*models.Snippet
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

	// Lopp through the pages filepath one by one
	for _, page := range pages{
		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
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