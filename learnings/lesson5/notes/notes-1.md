## The http.Handler interface

`handler` is an object which satisfies the
`http.Handler` interface:
```
type Handler interface {
 ServeHTTP(ResponseWriter, *Request)
}
```

To be a handler an object must have a ServeHTTP()
method:
```
ServeHTTP(http.ResponseWriter, *http.Request)
```

So in its simplest form a handler might look something like this:
```
type home struct {}
func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
 w.Write([]byte("This is my home page"))
}
```
Here we have an object (in this case it’s a home struct, but it could equally be a string or
function or anything else), and we’ve implemented a method with the signature
ServeHTTP(http.ResponseWriter, *http.Request) on it. That’s all we need to make a
handler.

You could then register this with a servemux using the Handle method like so:
```
mux := http.NewServeMux()
mux.Handle("/", &home{})
```
When this servemux receives a HTTP request for "/", it will then call the ServeHTTP() method
of the home struct — which in turn writes the HTTP response.


### Handler functions

In practice it’s far more common to write your handlers as a
normal function. Example:
```
func home(w http.ResponseWriter, r *http.Request) {
 w.Write([]byte("This is my home page"))
}
```
But this home function is just a normal function; it doesn’t have a `ServeHTTP()` method. So in
itself it isn’t a handler.

Instead we can transform it into a handler using the http.HandlerFunc() adapter.
```
mux := http.NewServeMux()
mux.Handle("/", http.HandlerFunc(home))
```
The `http.HandlerFunc()` adapter works by automatically adding a `ServeHTTP()` method to
the home function.


## Chaining handlers

The `http.ListenAndServe()` function takes a `http.Handler` object as the second
parameter but we’ve been passing in a servemux.
```
func ListenAndServe(addr string, handler Handler) error
```
Because the `servemux` also has a `ServeHTTP()` method, meaning that
it too satisfies the `http.Handler` interface.

The `servemux` is just being a special kind of handler,
which instead of providing a response itself passes the request on to a second handler. Such Chaining of handlers together is a very common idiom in Go.

In fact, what exactly is happening is this: When our server receives a new HTTP request, it calls
the servemux’s `ServeHTTP()` method. This looks up the relevant handler based on the
request URL path, and in turn calls that handler’s `ServeHTTP()` method. You can think of a Go
web application as a chain of `ServeHTTP()` methods being called one after another.


## Requests are handled concurrently

- All incoming HTTP requests are
served in their own goroutine.

- This make Go blazingly fast,
the downside you need to be aware of (and protect against) are race conditions when
accessing shared resources from your handlers.
