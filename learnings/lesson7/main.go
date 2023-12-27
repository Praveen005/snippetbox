package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct{
	infoLog *log.Logger
	errorLog *log.Logger
}

func main(){
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Initialize a new instance of application struct, containing the dependecies
	app := application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// swap the route declarations to use the application struct's methods as the handler functions

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}