package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet is the corresponding struct to database table snippet.
type Snippet struct {
    ID      int
    Title   string
    Content string
    Created time.Time
    Expires time.Time
}

// SnippetModelInterface defines the methods a SnippetModel struct should implement.
type SnippetModelInterface interface {
    Insert(title string, content string, expires int) (int, error)
    Get(id int) (Snippet, error)
    Latest(n int) ([]Snippet, error)
}

// SnippetModel wraps a sql.DB connection pool.
type SnippetModel struct {
    DB *sql.DB
}

// Insert inserts a new record in database table snippet.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
    stmt := `INSERT INTO snippet(title, content, created, expires) 
             VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

    result, err := m.DB.Exec(stmt, title, content, expires)
    if err != nil {
        return 0, err
    }

    // Use the LastInsertId() method on the result to get the ID of our newly inserted record.
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    // The ID returned has the type int64, so we convert it to an int type before returning.
    return int(id), nil
}

// Get returns a specific Snippet based on its ID.
func (m *SnippetModel) Get(id int) (Snippet, error) {
    stmt := `SELECT id, title, content, created, expires 
               FROM snippet 
              WHERE expires > UTC_TIMESTAMP() 
                AND id = ?`

    var s Snippet

    err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
    if err != nil {
        // If the query returns no rows, Scan() will return a sql.ErrNoRows error. We use the 
        // errors.Is() function to check for that error specifically, and return our own 
        // ErrNoRecord error instead.
        if errors.Is(err, sql.ErrNoRows) {
            return Snippet{}, ErrNoRecord
        } else {
            return Snippet{}, err
        }
    }

    return s, nil
}

// Latest returns n most recently created snippets.
func (m *SnippetModel) Latest(n int) ([]Snippet, error) {
    stmt := `SELECT id, title, content, created, expires 
               FROM snippet 
              WHERE expires > UTC_TIMESTAMP() 
              ORDER BY id DESC 
              LIMIT ?`

    rows, err := m.DB.Query(stmt, n)
    if err != nil {
        return nil, err
    }
    // We defer rows.Close() to ensure the sql.Rows resultset is always properly closed before this
    // method returns. This defer statement should come *after* you check for an error from the
    // Query() method. Otherwise, if Query() returns an error, you'll get a panic trying to close 
    // a nil resultset.
    defer rows.Close()

    var snippets []Snippet

    // Use rows.Next to iterate through the rows in the resultset. This prepares the first (and 
    // then each subsequent) row to be acted on by the rows.Scan() method. If iteration over all 
    // the rows completes, the resultset automatically closes itself and freesup the uderlying 
    // database connection.
    for rows.Next() {
        var s Snippet

        err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
        if err != nil {
            return nil, err
        }

        snippets = append(snippets, s)
    }

    // When the rows.Next() loop has finished we call rows.Err() to retrieve any error that was 
    // encountered during the iteration. It's important to call this - don't assume that a 
    // successful iteration was completed over the whole resultset.
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return snippets, nil
}