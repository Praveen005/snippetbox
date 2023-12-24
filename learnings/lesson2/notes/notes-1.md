## System-generated headers and content sniffing

- When sending a response Go will automatically set three system-generated headers for you:
Date and Content-Length and Content-Type.

- The Content-Type header is particularly interesting. Go will attempt to set the correct one for
you by content sniffing the response body with the http.DetectContentType() function. If
this function can’t guess the content type, Go will fall back to setting the header
Content-Type: application/octet-stream instead

- **The http.DetectContentType() function can’t distinguish JSON from plain text**.

- By default,
JSON responses will be sent with a ```Content-Type: text/plain;``` charset=utf-8 header

- We
can prevent this from happening by setting the correct header manually:

```
w.Header().Set("Content-Type", "application/json")
w.Write([]byte(`{"name":"Praveen"}`))
```

## Manipulating the header map

We can use ```w.Header().Set()/Add()/Del()/Get()/Vales()``` methods that you can use to read
and manipulate the header map

```
// Set a new cache-control header. If an existing "Cache-Control" header exists
// it will be overwritten.
w.Header().Set("Cache-Control", "public, max-age=31536000")


// In contrast, the Add() method appends a new "Cache-Control" header and can
// be called multiple times.
w.Header().Add("Cache-Control", "public")
w.Header().Add("Cache-Control", "max-age=31536000")


// Delete all values for the "Cache-Control" header.
w.Header().Del("Cache-Control")


// Retrieve the first value for the "Cache-Control" header.
w.Header().Get("Cache-Control")


// Retrieve a slice of all values for the "Cache-Control" header.
w.Header().Values("Cache-Control")

```


## Header canonicalization

Header canonicalization is a process that ensures uniformity in the representation of HTTP header names. When you’re using the Set(), Add(), Del(), Get() and Values() methods on the header map,
the header name will always be canonicalized using the
```textproto.CanonicalMIMEHeaderKey()``` function.

The canonicalization process follows these rules:

1. The first letter of the header name is converted to uppercase.

2. Any letter following a hyphen (-) is also converted to uppercase.

3. The rest of the letters are converted to lowercase.

This canonicalization process makes header names case-insensitive, which means that regardless of how you provide the header name (with different cases), it will be stored and treated in a standardized form.

```
// Adding headers using various cases
resp.Header.Set("content-type", "application/json")
resp.Header.Add("X-Custom-header", "value1")
```

It will be stored like:

```
Content-Type: application/json
X-Custom-Header: value1
```

There might be cases where you want to avoid this canonicalization and directly manipulate the underlying header map without any automatic capitalization. In such situations, you can access the raw header map, which has the type ```map[string][]string```, and modify it directly.
Example:

```
w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
```

> Note: If a ```HTTP/2``` connection is being used, Go will always automatically convert the
header names and values to lowercase for you as per the ```HTTP/2``` specifications.
>
> Example: if you set a header using ```resp.Header.Set("Content-Type", "application/json")``` and the request is sent over an ```HTTP/2``` connection, the actual headers sent in the ```HTTP/2``` frames will have the header field name ```"content-type"``` in lowercase.

### Suppressing system-generated headers

The ```Del()``` method doesn’t remove system-generated headers. To suppress these, you need
to access the underlying header map directly and set the value to nil. If you want to suppress
the Date header, for example, you need to write:

```
w.Header()["Date"] = nil
```