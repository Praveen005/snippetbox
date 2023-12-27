# Dependency Injection

There’s one more problem with our logging that we need to address. If we open up our
`handlers.go` file you’ll notice that the `home` handler function is still writing error messages
using Go’s standard logger, not the `errorLog` logger that we want to be using.

```
func home(w http.ResponseWriter, r *http.Request) {
    ...
    ts, err := template.ParseFiles(files...)

    if err != nil {
        log.Print(err.Error()) // This isn't using our new error logger.
        http.Error(w, "Internal Server Error", 500)
        return
    }

    err = ts.ExecuteTemplate(w, "base", nil)

    if err != nil {
        log.Print(err.Error()) // This isn't using our new error logger.
        http.Error(w, "Internal Server Error", 500)
    }
}
```

**How can we make our new `errorLog` logger available to our home
function from `main()`?**

Most web applications will have multiple dependencies
that their handlers need to access, such as a database connection pool, centralized error
handlers, and template caches.

**So, how can we make any
dependency available to our handlers?**

There are a few different ways to do this:

1. The simplest being is to just put the dependencies in
global variables.

2. But in general, it is good practice to inject dependencies into your handlers.
    - It makes your code more explicit, less error-prone and easier to unit test than if you use global
variables.

    - For applications where all your handlers are in the same package, like ours, a neat way to
inject dependencies is to put them into a custom application struct, and then define your
handler functions as methods against application.

 ```
// Define an application struct to hold the application-wide dependencies for the

// web application. For now we'll only include fields for the two custom loggers, but

// we'll add more to it as the build progresses.

type application struct {
    errorLog *log.Logger
    infoLog *log.Logger
}
```

-----------------
**`app.errorLog.Print(err.Error())`  vs  `app.errorLog.Print(err)`**

- `app.errorLog.Print(err.Error())` logs the string representation of the error obtained by calling `err.Error()`

- The `Error()` method of the error interface returns a string representing the error message.

- This is useful when you want to log the error message explicitly as a string.

- `app.errorLog.Print(err)` logs the error value directly, without converting it to a string
```
err := &CustomError{Code: 42, Message: "Something went wrong"}

app.errorLog.Print(err) // logs the entire CustomError value
```
