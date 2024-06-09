package models

import (
	"database/sql"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?

type Snippet struct{
	ID		int
	Title	string
	Content string
	Created time.Time		// once a time.Time value is created, its internal state cannot be changed.
	Expires time.Time
}


// Define SnippetModel which eraps the sql.DB connection pool
type SnippetModel struct{
	DB  *sql.DB
}


// This will insert a new snippet in the database and return the id of the snippet created
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`


	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// title, content and expiry values for the placeholder parameters. This
	// method returns a sql.Result type, which contains some basic
	// information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil{
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	// In Go, int and int64 are distinct types, and they are not interchangeable.
	// Many Go libraries and standard library functions expect int rather than int64.
	// Using int64 when int is expected would cause type mismatch 
	// result.LastInsertId() returns an int64 because database IDs can potentially be very large.
	return int(id), nil
}

// This will return a specific snippet based on id
func(m *SnippetModel) Get(id int)(*Snippet, error){
	return nil, nil
}