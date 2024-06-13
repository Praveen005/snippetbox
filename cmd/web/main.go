package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.01 (Project Setup and Creating
	// a Module) so that the import statement looks like this:
	// "{your-module-path}/internal/models". If you can't remember what module path you
	// used, you can find it at the top of the go.mod file.
	"github.com/Praveen005/snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// Add a new sessionManager field to the application struct.
type application struct{
	infoLog  		*log.Logger
	errorLog 		*log.Logger
	snippets 		*models.SnippetModel
	templateCache	map[string]*template.Template
	formDecoder 	*form.Decoder
	sessionManager  *scs.SessionManager
}

func main(){
	//flag.String function returns a pointer to a string
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	// Define a new command-line flag for the MySQL DSN(data source name: depend on which database and driver you’re using.) string.
	dsn := flag.String("dsn", "web:p123@/snippetbox?parseTime=true", "MYSQL data source name")
	flag.Parse()	


	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil{
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits
	defer db.Close()


	// Initializa a new template cache
	templateCache, err  := newTemplateCache()
	if err != nil{
		errorLog.Fatal(err)
	}


	// Initialize a decoder instance
	formDecoder := form.NewDecoder()


	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions automatically expire 12 hours
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour


	// And add the session manager to our application dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)

	if err != nil{
		return nil, err
	}

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return db, nil
}