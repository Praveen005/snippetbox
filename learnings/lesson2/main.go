package main

import (
	"log"
	"net/http"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body
// Check if the current request URL path exactly matches "/". If it doesn't, use
 // the http.NotFound() function to send a 404 response to the client.
 // Importantly, we then return from the handler. If we don't return the handler
 // would keep executing and also write the "Hello from SnippetBox" message.
func home(w http.ResponseWriter, r * http.Request){
	if r.URL.Path != "/"{
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r * http.Request){
	w.Write([]byte("Display a specific snippet.."))
}

// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r * http.Request){

	// If it's not, use the w.WriteHeader() method to send a 405 status
	// code and the w.Write() method to write a "Method Not Allowed"
	// response body. We then return from the function so that the
	// subsequent code is not executed.

	if r.Method != "POST"{
		// we can use the constant http.MethodPost instead of the string "POST", and the constant http.StatusMethodNotAllowed instead of the integer 405.
		// w.Header().Set("Allow", "POST")
		w.Header().Set("Allow", http.MethodPost)


		// In practice, it’s quite rare to use the w.Write() and w.WriteHeader() methods
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		// so, we use make use of another function http.error() to do that for me 

		
		// http.Error(w, "Method No Allowed", 405)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	w.Write([]byte("Create a new snippet.."))
}

func main(){

	// Use the http.NewServeMux() function to initialize a new servemux, then
	//  register the home function as the handler for the "/" URL pattern.

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)


	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)

}