##### Creating a module

you can think of a module path as basically being a canonical name or identifier for your project.

You can pick almost any string as your module path, but the important thing to focus on is
uniqueness. 

To avoid potential import conflicts with other people’s projects or the standard library in the future, you want to pick a module path that is globally unique and unlikely to be used by anything else.

In my case, a clear, succinct and unlikely-to-be-used-by-anything-else module path for this
project would be snippetbox.praveen, and I’ll use this throughout the rest of the book.

```go mod init snippetbox.praveen```

or can simply have:

```go mo init snippetbox```

Additional Info:

    If you’re creating a project which can be downloaded and used by other people and
    programs, then it’s good practice for your module path to equal the location that the code
    can be downloaded from.

    For instance, if your package is hosted at https://github.com/Praveen005/snippetbox then the module path
    for the project should be github.com/Praveen005/snippetbox.

