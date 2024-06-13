# Stateful HTTP

Ref: pg. 232:

PS: We are using [alexedwards/scs](https://github.com/alexedwards/scs) for session management.

```
app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
```
### Put Method:

- The `Put` method is used to add or update a value in the session data. It takes three parameters:

    1. `r.Context()`: The context of the current HTTP request. The session manager uses this context to associate the session data with the current user's session.

    2. `"flash"`: The key under which the value is stored. In this case, "flash" is the key.

    3. `"Snippet successfully created!"`: The value to be stored in the session. This is a message indicating that a snippet was successfully created.

### Purpose of the Code:

The specific line of code is setting a `"flash"` message in the session. A flash message is a temporary message that is typically used to display success, error, or informational messages to users after they perform some action (such as submitting a form).

> The flash message is removed from the session storage after it has been retrieved and displayed to the user. This ensures that the flash message is displayed only once and does not persist across multiple requests.

### PopString() method:

Using this method, on the redirected page, the flash message is retrieved from the session, displayed to the user, and then removed from the session storage to ensure it is only shown once.


### Session Expiry:

Session expiry usually means, session is invalid and should not be used for authentication or storing temporary data anymore.

The deletion of session data from the database might not happen right at the moment of expiry. Instead, may get cleaned up periodically by the session management system to ensure efficient use of database resources.

### Conlusion:

So, what happens in our application is that the LoadAndSave() middleware checks each
incoming request for a session cookie. If a session cookie is present, it reads the session token
and retrieves the corresponding session data from the database (while also checking that the
session hasn’t expired). It then adds the session data to the request context so it can be used
in your handlers.
