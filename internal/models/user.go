package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// User is the corresponding struct to database table user.
type User struct {
    ID             int
    Name           string
    Email          string
    HashedPassword string
    Created        time.Time
}

// UserModelInterface defines the methods a UserModel struct should implement.
type UserModelInterface interface {
    Insert(name, email, password string) error
    Get(id int) (User, error)
    Exists(id int) (bool, error)
    Authenticate(email, password string) (int, error)
    UpdatePassword(id int, currentPassword, newPassword string) error
}

// UserModel wraps a *sql.DB connection pool.
type UserModel struct {
    DB *sql.DB
}

// Insert inserts a record in the user table.
func (m *UserModel) Insert(name, email, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

    stmt := `INSERT INTO user(name, email, hashed_password, created) 
             VALUES (?, ?, ?, UTC_TIMESTAMP())`

    _, err = m.DB.Exec(stmt, name, email, hashedPassword)
    if err != nil {
        // If this returns an error, we use the errors.As() function to check whether the rror has
        // the type *mysql.MySQLError. If it does, the error will be assigned to the mySQLError
        // variable. We can check whether or not the error relates to our uc_user_email constraint
        // by checking if the error code equals 1062 (ER_DUP_ENTRY) and the contents of the error
        // message string. If it does, we return an ErrDuplicateEmail error.
        var mySQLError *mysql.MySQLError
        if errors.As(err, &mySQLError) {
            if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "uc_user_email") {
                return ErrDuplicateEmail
            }
        }

        return err
    }

    return nil
}

// Get returns a specific User based on its ID.
func (m *UserModel) Get(id int) (User, error) {
    stmt := `SELECT id, name, email, hashed_password, created 
               FROM user
              WHERE id = ?`

    var u User

    err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Created)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return User{}, ErrNoRecord
        } else {
            return User{}, err
        }
    }

    return u, nil
}

// Exists checks if a user exists based on its ID.
func (m *UserModel) Exists(id int) (bool, error) {
    stmt := `SELECT EXISTS(SELECT true FROM user WHERE id = ?)`

    var exists bool

    err := m.DB.QueryRow(stmt, id).Scan(&exists)

    return exists, err
}

// Authenticate verifies whether a user exists based on the provided email and password.
// It returns the relevant user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
    stmt := `SELECT id, hashed_password
               FROM user 
              WHERE email = ?`

    var (
        id             int
        hashedPassword string
    )

    err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return 0, ErrInvalidCredentials
        } else {
            return 0, err
        }
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return 0, ErrInvalidCredentials
        } else {
            return 0, err
        }
    }

    return id, nil
}

// UpdatePassword updates a user's password.
func (m *UserModel) UpdatePassword(id int, currentPassword, newPassword string) error {
    stmt := `SELECT hashed_password 
               FROM user 
              WHERE id = ?`

    var currentHashedPassword string

    err := m.DB.QueryRow(stmt, id).Scan(&currentHashedPassword)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return ErrNoRecord
        } else {
            return err
        }
    }

    err = bcrypt.CompareHashAndPassword([]byte(currentHashedPassword), []byte(currentPassword))
    if err != nil {
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return ErrInvalidCredentials
        } else {
            return err
        }
    }

    stmt = `UPDATE user 
            SET hashed_password = ? 
            WHERE id = ?`

    newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
    if err != nil {
        return err
    }

    _, err = m.DB.Exec(stmt, string(newHashedPassword), id)

    return err
}
