# Isolating the application routes

Our `main()` function is beginning to get a bit crowded, so to keep it clear and focused I’d like
to move the route declarations for the application into a standalone `routes.go` file.


The responsibilities of our main() function will now be limited to:

- Parsing the runtime configuration settings for the application.

- Establishing the dependencies for the handlers.

- Running the HTTP server.