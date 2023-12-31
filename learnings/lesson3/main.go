//let’s update the snippetView handler so that it accepts an id query string parameter from the user

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r * http.Request){
	if r.URL.Path != "/"{
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Welcome to Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request){
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404 page
	// not found response.

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	
	if err != nil  || id < 1{
		http.NotFound(w, r)
		return
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Display specific snippet with ID %d..", id);
}

func snippetCreate(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		w.Header().Set("Allow", http.MethodPost)

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet.."))
}

func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)

	if err != nil{
		log.Fatal(err)
	}
}