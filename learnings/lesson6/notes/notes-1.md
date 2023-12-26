# Configuration and error handling

Here we:

- Set configuration settings for our application at runtime in an easy and idiomatic way
using command-line flags.

- Improve our application log messages to include more information, and manage them
differently depending on the type (or level) of log message.

- Make dependencies available to our handlers in a way that’s extensible, type-safe, and
doesn’t get in the way when it comes to writing tests.

- Centralize error handling so that we don’t need to repeat ourself when writing code.

Our web application’s main.go file currently contains a couple of hard-coded configuration
settings:

- The network address for the server to listen on (currently ":4000")

- The file path for the static files directory (currently "./ui/static")

> Having these hard-coded isn’t ideal. There’s no separation between our configuration
settings and code, and we can’t change the settings at runtime (which is important if you
need different settings for development, testing and production environments)


## Command-line flags

In Go, a common and idiomatic way to manage configuration settings is to use command-line
flags when starting an application. For example:
```
go run ./cm/web -addr=":80"
```

The easiest way to accept and parse a command-line flag from an application is like this:
```
addr := flag.String("addr", ":4000", "HTTP network address")
```
- This essentially defines a new command-line flag with the name addr, a default value of
":4000" and some short help text explaining what the flag controls. The value of the flag will
be stored in the addr variable at runtime.

>Note: Ports 0-1023 are restricted and (typically) can only be used by services which have
root privileges. If you try to use one of these ports you should get a
bind: permission denied error message on start-up.

- Command-line flags are completely optional. If we run the application with no
`-addr` flag the server will fall back to listening on address `:4000` (which is the default value we
specified).

**Automated help**

- We can use the -help flag to list all the available command-line
flags for an application and their accompanying help text.
Example:
```
$ go run ./cmd/web -help
Usage of C:\Users\PRAVEE~1\AppData\Local\Temp\go-build2456658510\b001\exe\web.exe:
  -addr string
        Http Network Address (default ":4000")
```

## Environment variables

If we want, we can store your configuration settings in environment variables and access
them directly from our application by using the `os.Getenv()` function, like:

```
addr := os.Getenv("SNIPPETBOX_ADDR")
```
But this has some drawbacks compared to using command-line flags. You can’t specify a
default setting (the return value from `os.Getenv()` is the empty string if the environment
variable doesn’t exist), you don’t get the `-help` functionality that you do with command-line
flags, and the return value from `os.Getenv()` is always a string — you don’t get automatic
type conversions like you do with `flag.Int()` and the other command line flag functions.


We can get the best of both worlds by passing the environment variable as a
command-line flag when starting the application. For example:
```
$ export SNIPPETBOX_ADDR=":9999"
$ go run ./cmd/web -addr=$SNIPPETBOX_ADDR
2022/01/29 15:54:29 Starting server on :9999
```

## Boolean flags
For flags defined with `flag.Bool()` omitting a value is the same as writing `-flag=true`. The
following two commands are equivalent:
```
$ go run ./example -flag=true
$ go run ./example -flag
```

We must explicitly use `-flag=false` if you want to set a boolean flag value to false.


## Pre-existing variables

 In Go, you can use the flag package to parse command-line arguments and store values directly into pre-existing variables. This is achieved using functions like `flag.StringVar()`, `flag.IntVar()`, and `flag.BoolVar()`. 

 ```
 type config struct {
  addr string
  staticDir string
}

...

var cfg config

flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

flag.Parse()
 ```



## Leveled logging

- Till now from our main.go file we’re outputting log messages using the `log.Printf()` and
`log.Fatal()` functions.

- Both these functions output messages via Go’s standard logger, which — by default —
prefixes messages with the local date and time and writes them to the standard error stream
(which displays in our terminal window).

- **The log.Fatal() function will also call
os.Exit(1) after writing the message, causing the application to immediately exit.**

- we can break apart our log messages into two distinct types — or levels:
  - Informational message
  - Error message
  ```
  log.Printf("Starting server on %s", *addr) // Information message
  err := http.ListenAndServe(*addr, mux)
  log.Fatal(err) // Error message
  ```


Let’s improve our application by adding some leveled logging capability:
  - We will prefix informational messages with "INFO" and output the message to standard
out (stdout).

  - We will prefix error messages with "ERROR" and output them to standard error (stderr),
along with the relevant `file name` and `line number` that called the logger (to help with
debugging).

There are a couple of different ways to do this, but a simple and clear approach is to use the
`log.New()` function to create two new custom loggers.

```
infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
```