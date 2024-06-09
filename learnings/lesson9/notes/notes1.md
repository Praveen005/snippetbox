# Database-driven responses

- Correctly using placeholder parameters can prevent SQL Injection
- A database driver acts as a ‘middleman’ between MySQL and your Go application


### Collation
Collation in the context of databases refers to a set of rules that determine how string data is sorted and compared.

Common Collations
- utf8mb4_unicode_ci: A widely-used collation that is case-insensitive, accent-insensitive(characters like é, è, and e are treated as equivalent), and follows the Unicode standard. Suitable for most applications that need to support multiple languages.

- utf8mb4_general_ci: Similar to utf8mb4_unicode_ci, but with slightly less accurate sorting rules for certain characters. Faster but less precise.

- utf8mb4_bin: A binary collation that treats strings as sequences of bytes, making it case-sensitive and accent-sensitive.

Example Usage:
```
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```
## 
Table we created:
```
+---------+--------------+------+-----+---------+----------------+
| Field   | Type         | Null | Key | Default | Extra          |
+---------+--------------+------+-----+---------+----------------+
| id      | int          | NO   | PRI | NULL    | auto_increment |
| title   | varchar(100) | NO   |     | NULL    |                |
| content | text         | NO   |     | NULL    |                |
| created | datetime     | NO   | MUL | NULL    |                |
| expires | datetime     | NO   |     | NULL    |                |
+---------+--------------+------+-----+---------+----------------+
```
What does `MUL` in `Key` column mean?

"MUL" stands for "multiple".

> In the context of a database index, let's say you have a column named "City" in a table that stores customer data. Each row in this table represents a different customer, and the "City" column indicates the city where each customer lives.
>
> Now, imagine you have several customers who all live in the same city, for example, New York. If you create an index on the "City" column to improve search performance when querying based on city names, MySQL will create an index entry for each occurrence of "New York" in the "City" column.
>
> In this scenario:
- The index is on the "City" column.
- The value "New York" appears multiple times in the "City" column.
- Each occurrence of "New York" will have its own entry in the index.
>
> When the "Key" column in the index information displays "MUL", it indicates that the index is on a column where multiple occurrences of the same value are possible. In this case, "MUL" stands for "multiple", meaning multiple rows can have the same value in the indexed column.


##

ref: page 105:

- In Go, `int` and `int64` are distinct types, and they are not interchangeable.
- `int` is a platform-dependent type. On a 32-bit system, `int` is 32 bits wide. On a 64-bit system, it is 64 bits wide.
- Many Go libraries and standard library functions expect `int` rather than `int64`.
- Using `int64` when `int` is expected would cause type mismatch errors.
- `result.LastInsertId()` returns an `int64` because database IDs can potentially be very large.



## What is `http.StatusSeeOther`?

ref: page: 106

`http.StatusSeeOther` is an HTTP status code constant in Go's `net/http` package. It corresponds to the HTTP status code 303 See Other.

### HTTP Status Code 303 See Other

- **Status Code**: 303
- **Name**: See Other
- **Description**: The response to the request can be found under a different URI, and should be retrieved using a GET method on that resource.
- **Use Case**: It is typically used in conjunction with the `Location` header in a response to indicate that the requested resource has been replaced with a different resource and the client should retrieve it using a new URI provided in the response.

### Purpose in `http.Redirect`

In the context of `http.Redirect`, `http.StatusSeeOther` is used to instruct the client to redirect to a different URI. When the client receives a response with this status code and a `Location` header, it should issue a GET request to the URI specified in the `Location` header to retrieve the resource.


### Use Case Example

After a user submits a form to create a new snippet, and upon successful creation, you want to redirect the user to the newly created snippet's page. You can use `http.StatusSeeOther` to redirect the user's browser to the page displaying the newly created snippet.

```go
// After creating the snippet, redirect the user to the snippet's page
http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
```

This results in the browser issuing a GET request to the URI `/snippet/view?id=<id>` to retrieve and display the newly created snippet.


## Creating a new user in mysql

1. Open mysql terminal
wsl:
```
sudo mysql -u root -p
```
Poweshell:
```
mysql -u praveen -p
```
2. Enter your password.

3. Create a new user:
```
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';
```

If you made any error, and want to drop this new user and create a new one:
1. If you are in mysql terminal, exit by:
```
exit
```
2. Now navigate to root user which has root privileges:
```
sudo mysql -u root -p
```
3. Enter your password

4. Get the list of users, and find yours:
```
SELECT user, host FROM mysql.user;
```
5. Now drop the user you want:
```
DROP USER 'user1'@'localhost';
```
6. Create a new one again now. yayy!

### To clear your mysql terminal:
wsl:     `ctrl + L`

Windows: `system cls;`