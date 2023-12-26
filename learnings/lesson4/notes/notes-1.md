## HTML templating

> Note: The `.tmpl` extension doesn’t convey any special meaning or behavior here. It has
only chosen because it’s a nice way of making it clear that the file
contains a Go template when you’re browsing a list of files.
>
> We can use the extension `.html` instead

Go’s `html/template` package, which provides a family of functions for
safely parsing and rendering HTML templates.

## `.Execute()` method in the `html/template` package 

Purpose:

- Executes (renders) an HTML template, combining a template structure with dynamic data to produce the final HTML output.

- Essential for generating dynamic web content in Go web applications.

Syntax:

```
func (t *Template) Execute(wr io.Writer, data interface{}) error
```

Parameters:

- `wr`: An `io.Writer` interface representing the destination where the rendered HTML output will be written. It can be a web response, a file, a buffer, or any other object that accepts written data.

- `data`: An interface{} value containing the data to be injected into the template during rendering. It's typically a struct, a map, or a slice that holds the values to be used within the template.

Return Values:

- `error`: An error value, if any, encountered during the execution process.

Key Points:

- Template Structure: The template must be parsed beforehand using `template.ParseFiles()` or `template.New().Parse()`.

- Data Injection: Template actions (like `{{.}}`, `{{range}}`, etc.) within the template extract and process data from the provided data value.

Common Template Actions:

- `{{.}}`: Prints the value of the current field or item in a loop.

- `{{range}}`: Iterates over a collection (e.g., slice or map) and executes a block of template code for each item.

- `{{if}}` and `{{else}}`: Conditionally executes blocks of code based on boolean expressions.

- With Pipelines: Chains multiple actions together to transform or format data before insertion.

Example:
```
<h1>Hello, {{.Name}}!</h1>  
<ul>
    {{range .Items}}  <li>{{.}}</li>  {{end}}
</ul>
```

As we add more pages to this web application there will be some shared, boilerplate, HTML
markup that we want to include on every page.

To save us typing and prevent duplication, it’s a good idea to create a base (or master)
template which contains this shared content.

```
{{define "base"}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>{{template "title" .}} - Snippetbox</title>
  </head>
  <body>
    <header>
      <h1><a href="/">Snippetbox</a></h1>
    </header>
    <main>{{template "main" .}}</main>
    <footer>Powered by <a href="https://golang.org/">Go</a></footer>
  </body>
</html>
{{end}}
```

- `{{define "base"}}...{{end}}` action to define a distinct named
template called base, which contains the content we want to appear on every page.

- `{{template "title" .}}` and `{{template "main" .}}` actions to
denote that we want to invoke other named templates (called title and main) at a particular
point in the HTML.

> Note: The dot at the end of the `{{template "title" .}}` action
represents any dynamic data that one wants to pass on the fly.

We define `title` and `main` as follows:

```
{{define "title"}}Home{{end}}


{{define "main"}}
 <h2>Latest Snippets</h2>
 <p>There's nothing to see here yet!</p>
{{end}}
```


For some applications you might want to break out certain bits of HTML into partials that can
be reused in different pages or layouts.

To illustrate, let’s create a partial containing the
primary navigation bar for our web application:

`File: ui/html/partials/nav.tmpl`
```
{{define "nav"}}
 <nav>
 <a href='/'>Home</a>
</nav>
{{end}}
```


Go also provides a {{block}}...{{end}} action

```
{{define "base"}}
    <h1>An example template</h1>
    {{block "sidebar" .}}
        <p>My default sidebar content</p>
    {{end}}
{{end}}
```
Te content between  `{{block}}` and `{{end}}` actions represent the default template that will be invoked, if the user specified templete somehow don't get rendered or is missing.

So, if you don't define the `"sidebar"` templete somewhere(like, `{{define "sidebar"}}Sidebar hun mai!{{end}}`), the default one between `{{block}}` and `{{end}}` will be rendered oherwise, the defined one will be shown.


## The http.Fileserver handler

Go’s net/http package ships with a built-in http.FileServer handler which you can use to
serve files over HTTP from a specific directory

### http.StripPrefix
```
fileServer := http.FileServer(http.Dir("./ui/static/"))
http.Handle("/static/", http.StripPrefix("/static", fileServer))
```
- Static files will be searched in ./ui/static directory
- our URL route is /static/ and it will look for file /static/css/style.css in ./ui/static
- so new path will become, /static/static/css/style.css, which doesn't exist
- so we strip /static prefix from the URL, and it becomes /static/css/style.css, which is correct
- means it is looking for /css/style.css in ./ui/staic/
- This is one reason why we need it.

## Features and functions of Go's file server

- It sanitizes all request paths by running them through the `path.Clean()` function before
searching for a file. This removes any . and .. elements from the URL path, which helps to
stop directory traversal attacks.

    For example, they might submit a path like `"../../sensitive_file.txt"` to go up two directory levels.

    if a request includes `"/../../some-sensitive-file"`, `path.Clean()` would sanitize it to `"/some-sensitive-file"`, avoiding unauthorized access.

- Range requests are fully supported. This is great if your application is serving large files
and you want to support resumable downloads.

    If a download is interrupted or fails, the client can resume the download by requesting the missing or incomplete range. This is especially useful for large files where downloading from the beginning would be inefficient.

    Bandwidth Efficiency: Range requests allow clients to download specific parts of a file, reducing the amount of data transferred and improving bandwidth efficiency.

    User Experience: Users can start interacting with the content before the entire resource is downloaded, leading to a better user experience, especially for multimedia content.

    > Note: The 206 Partial Content status code indicates that the server is responding with only a portion of the requested resource.


- The `Last-Modified` and `If-Modified-Since` headers are transparently supported. If a file
hasn’t changed since the user last requested it, then `http.FileServer` will send a
`304 Not Modified` status code instead of the file itself. This helps reduce latency and
processing overhead for both the client and server


- The Content-Type is automatically set from the file extension using the
mime.TypeByExtension() function. You can add your own custom extensions and content
types using the mime.AddExtensionType() function if necessary.
