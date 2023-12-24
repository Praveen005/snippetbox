## Project Structure:
```
.
└── Snippetbox/
    ├── cmd/
    │   └── web/
    │       ├── handler.go
    │       └── main.go
    ├── Internal
    ├── learnings
    ├── ui/
    │   ├── html
    │   └── static
    └── go.mod

```
> PS: I'm using Powershell

```
New-Item -Path "cmd/web", "internal", "ui/html", "ui/static" -ItemType Directory -Force
```
```
New-Item -Path "cmd/web/main.go" -ItemType File
```
```
New-Item -Path "cmd/web/handlers.go" -ItemType File 
```

elsewhere you can run,
```
$ cd $HOME/code/snippetbox
$ mkdir -p cmd/web internal ui/html ui/static
$ touch cmd/web/main.go
$ touch cmd/web/handlers.go
```

- The `cmd` directory will contain the _application-specific_ code.

- The `internal` directory will contain the ancillary _non-application-specific_ code used in the
project. We’ll use it to hold potentially reusable code like validation helpers and the SQL
database models for the project.

- The `ui` directory will contain the user-interface assets used by the web application.
Specifically, the `ui/html` directory will contain HTML templates, and the `ui/static`
directory will contain static files (like CSS and images).


**Benefits:**

- It gives a clean separation between Go and non-Go assets.

- All the Go code we write will
live exclusively under the cmd and internal directories, leaving the project root free to
hold non-Go assets like UI files, makefiles and module definitions (including our go.mod
file).

- It scales really nicely if you want to add another executable application to your project.
For example, you might want to add a CLI (Command Line Interface) to automate some
administrative tasks in the future. With this structure, you could create this CLI application
under cmd/cli and it will be able to import and reuse all the code you’ve written under
the internal directory.

### The internal directory:

- The directory name `internal` carries a special meaning and
behavior in Go: any packages which live under this directory can only be imported by code
inside the parent of the `internal` directory.

- In our case, this means that any packages which
live in `internal` can only be imported by code inside our snippetbox project directory.

- This means that any packages under `internal` cannot be
imported by code outside of our project.

- This is useful because it prevents other codebases from importing and relying on the
(potentially unversioned and unsupported) packages in our internal directory — even if the
project code is publicly available somewhere like GitHub.