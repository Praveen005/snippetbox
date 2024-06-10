package main


import(
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


		/// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		// parse the files into template set
		ts, err := template.ParseFiles(files...)
		if err != nil{
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		//  (like 'home.tmpl') as the key.
		cache[name] = ts
	}
	// return the map
	return cache, nil
}