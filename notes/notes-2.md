## Web application basics

Three absolute essentials for our web applications are:

1. Handler: 
    - If one is coming from an MVC-background, we can
think of handlers as being a bit like controllers.

    - They’re responsible for executing our
application logic and for writing HTTP response headers and bodies.

2. Router (or servemux in Go terminology):
    - This stores a mapping
between the URL patterns for your application and the corresponding handlers.

   - Usually
you have one servemux for your application containing all your routes.

3. Web Server:
    - One of the great things about Go is that you can
establish a web server and listen for incoming requests as part of your application itself.
You don’t need an external third-party server like Nginx or Apache.



### [log Package](https://pkg.go.dev/log):

- Package log implements a simple logging package. 

- It defines a type, Logger, with methods for formatting output. 

- It also has a predefined 'standard' Logger accessible through helper functions Print[f|ln], Fatal[f|ln], and Panic[f|ln], which are easier to use than creating a Logger manually. 

- That logger writes to standard error and prints the date and time of each logged message. 

- Every log message is output on a separate line: if the message being printed does not end in a newline, the logger will add one. 

- The Fatal functions call os.Exit(1) after writing the log message. 

- The Panic functions call panic after writing the log message.


> Note: The home handler function is just a regular Go function with two parameters. The
http.ResponseWriter parameter provides methods for assembling a HTTP response
and sending it to the user, and the *http.Request parameter is a pointer to a struct
which holds information about the current request (like the HTTP method and the URL
being requested). 


Each time the server receives a new HTTP request it will pass the request on to the
servemux and — in turn — the servemux will check the URL path and dispatch the request to
the matching handler.

Go’s servemux treats the URL
pattern "/" like a catch-all. So at the moment all HTTP requests to our server will be
handled by the home function, regardless of their URL path.


The TCP network address that you pass to http.ListenAndServe() should be in the format
"host:port". If you omit the host (like with ":4000") then the server will listen on all
our computer’s available network interfaces.


The following commands are all equivalent:
```github.com/Praveen005/snippetbox``` here is the module path
```
go run .
go run main.go
go run github.com/Praveen005/snippetbox
```


