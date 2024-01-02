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
	//flag.String function returns a pointer to a string
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Initialize a new instance of application struct, containing the dependecies
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}
	/*
		The http.Server type represents an HTTP server. It has fields such as Addr (the address to listen on), Handler (the handler to invoke), and other configuration options like:
		- ErrorLog(An optional logger for errors
		- ReadTimeout & WriteTimeout: Timeout values for reading and writing requests
		- MaxHeaderBytes: The maximum allowed size of the request headers
	*/

	//we use the http.Server struct literal to create a new server object
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)

	// The http.ListenAndServe function starts an HTTP server
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}