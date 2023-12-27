package main

import(
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"path"
)
// Change the signature of the home handler so it is defined as a method against
// *application.

func (app *application) home(w http.ResponseWriter, r * http.Request){
	if r.URL.Path != "/"{
		http.NotFound(w, r)
		return
	}
	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	files := []string{
		"./ui/html/pages/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",

	}

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file
	// paths as a variadic parameter?
	ts, err := template.ParseFiles(files...)

	if err != nil {
		//Because the home handler function is now a method against application
		//It can access its fields, including the error logger. We'll write the log message to this instead of standard logger.
		app.errorLog.Print(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base" 
	// template as the response body.
	/*
		So now, instead of containing HTML directly, our template set contains 3 named templates —
		base, title and main. We use the ExecuteTemplate() method to tell Go that we specifically
		want to respond using the content of the base template (which in turn invokes our title and
		main templates).
	*/
	err = ts.ExecuteTemplate(w, "base", nil)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

func snippetView(w http.ResponseWriter, r * http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1{
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r * http.Request){
	if r.Method != http.MethodPost{
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}

func downloadHandler(w http.ResponseWriter, r* http.Request){
	path := path.Clean("./notes-1.md")

	// w.Header().Set("Content-Disposition", "inline; filename=notes1.md") //displays the file in browser
	w.Header().Set("Content-Disposition", "attachment; filename=notes1.md") //downloads the file
	http.ServeFile(w, r, path)
}