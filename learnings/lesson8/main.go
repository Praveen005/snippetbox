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
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}