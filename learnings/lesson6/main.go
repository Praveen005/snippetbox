package main

import (
	"flag"
	"log"
	"net/http"
)

func main(){
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "Http Network Address")
	
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()


	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	mux.HandleFunc("/download", downloadHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting server on %s", *addr)

	err := http.ListenAndServe(*addr, mux)

	log.Fatal(err)
}