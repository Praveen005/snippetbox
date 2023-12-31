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
Run server in one terminal and try running it again in another, you will get the following:

![Info & Error](./assets/logs.png)


> Tip: If you want to include the full file path in your log output, instead of just the file
name, you can use the `log.Llongfile` flag instead of `log.Lshortfile` when creating
your custom logger. You can also force your logger to use UTC datetimes (instead of
local ones) by adding the `log.LUTC` flag.


## Decoupled logging

- A big benefit of logging your messages to the standard streams (stdout and stderr)
is that our application and logging are decoupled.

- Our application itself isn’t concerned
with the routing or storage of the logs, and that can make it easier to manage the logs
differently depending on the environment.

- During development, it’s easy to view the log output because the standard streams are
displayed in the terminal.

- In staging or production environments, you can redirect the streams to a final destination for
viewing and archival.
  - This destination could be on-disk files, or a logging service such as
Splunk.
  - Either way, the final destination of the logs can be managed by our execution
environment independently of the application.
  -  we can redirect the `stdout` and `stderr` streams to on-disk files when starting
the application like:
  ```
  go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
  ```

  > Note: Using the double arrow >> will append to an existing file, instead of truncating it
when starting the application,  meaning it adds the new output to the end of the existing file content.

In Unix-like systems, file descriptors are used to represent open files or I/O streams. 
- 0 - Standard Input (stdin): This is where a program receives input by default.

- 1 - Standard Output (stdout): This is where a program sends its normal output by default. Example: `>>/tmp/info.log`

- 2 - Standard Error (stderr): This is where a program sends error messages and diagnostic output by default. Example: `2>>/tmp/error.log`

If you're faced with Error: `out-file : Could not find a part of the path`
make the directory by:
```
mkdir tmp
```
Then try:
```
go run ./cmd/web >>./tmp/info.log 2>>./tmp/error.log
```



## The http.Server error log

- By default, if Go’s HTTP server
encounters an error it will log it using the standard logger.

- For consistency it’d be better to
use our new errorLog logger instead.

- To make this happen we need to initialize a new http.Server struct containing the
configuration settings for our server, instead of using the http.ListenAndServe() shortcut.

```
srv := &http.Server{
  Addr: *addr,
  ErrorLog: errorLog,
  Handler: mux,
}

err := srv.ListenAndServe()
errorLog.Fatal(err)
```

## Additional logging methods

As a rule of thumb, you should avoid using the `Panic()` and `Fatal()` variations outside of
your `main()` function — it’s good practice to return errors instead, and only panic or exit
directly from `main()`.


Custom loggers created by `log.New()` are concurrency-safe. You can share a single logger and
use it across multiple goroutines and in your handlers without needing to worry about race
conditions.
```
customLogger := log.New(os.Stdout, "CUSTOM: ", log.Ldate|log.Ltime)

go customLogger.Print("Hello Ji, Ram Ram!")
go customLogger.Print("Hello Bhai, Ram Ram!")
go customLogger.Print("Hello Didi, Ram Ram!")
// Remember, in a real-world scenario, you'd use synchronization mechanisms to wait for goroutines to finish
```

If we have multiple loggers writing to the same destination then we need to be
careful and ensure that the destination’s underlying Write() method is also safe for
concurrent use. Can be done using:

1. Synchronize Access:
    - Use synchronization mechanisms to ensure that only one goroutine at a time is writing to the destination. This can be achieved using a mutex (from the sync package) to lock access during writes.

2. Using a buffered channel can provide a safe way to handle concurrent writes.

### Logging to a file:

It is recommended to log the output to standard streams and
redirect the output to a file at runtime. But if one doesn’t want to do this, we can always open a
file in Go and use it as our log destination. Example:

First create the directory `logs` then use the following:

```
f, err := os.OpenFile("./logs/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
if err != nil {
 log.Fatal(err)
}
defer f.Close()
infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
```

Result:
![logs written a file](./assets/logs_in_file.png)

**About the above snippet:**

Opening or Creating a Log File:

  - It attempts to open or create a file named `./logs/info.log` using the `os.OpenFile` function.

  - The `os.O_RDWR` flag indicates that the file should be opened for both reading and writing.

  - The `os.O_CREATE` flag indicates that the file should be created if it does not exist.

  - By adding `os.O_APPEND` to the flags, the file will be opened in append mode, allowing new log entries to be appended to the existing content instead of overwriting it.

  - The file is created with read and write permissions `(0666)`.

Checking for Errors:

  - If there is an error opening or creating the file, it logs the error using log.Fatal and terminates the program.

  - The `log.Fatal` function prints the error and then calls `os.Exit(1)`, causing the program to exit with a non-zero status code.

Deferring File Closure:

  - Regardless of whether an error occurred or not, the `defer` statement ensures that the file is closed when the surrounding function (`main` or another function that encapsulates this code) exits.