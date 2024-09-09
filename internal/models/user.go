package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// User is the corresponding struct to the "user" table.
type User struct {
    ID             int
    Name           string
    Email          string
    HashedPassword []byte
    Created        time.Time
}

type UserModelInterface interface {
    Insert(name, email, password string) error
    Authenticate(email, password string) (int, error)
    Exists(id int) (bool, error)
    Get(id int) (User, error)
    UpdatePassword(id int, currentPassword, newPassword string) error
}

// UserModel wraps a database connection pool.
type UserModel struct {
    DB *sql.DB
}

// Insert inserts a record to the "user" table.
func (m *UserModel) Insert(name, email, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

    stmt := `INSERT INTO user (name, email, hashed_password, created) 
             VALUES(?, ?, ?, UTC_TIMESTAMP())`

    _, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
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

// Authenticate verifies whether a user exists with the provided email and password. 
// This will return the relevant user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
    var id int
    var hashedPassword []byte

    stmt := "SELECT id, hashed_password FROM user WHERE email = ?"

    err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return 0, ErrInvalidCredentials
        } else {
            return 0, err
        }
    }

    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    if err != nil {
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return 0, ErrInvalidCredentials
        } else {
            return 0, err
        }
    }

    return id, nil
}

// Exists checks if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
    var exists bool

    stmt := "SELECT EXISTS(SELECT true FROM user WHERE id = ?)"

    err := m.DB.QueryRow(stmt, id).Scan(&exists)

    return exists, err
}

// Get returns a User struct containing all the information for this user except for their hashed 
// password. If no user is found with the ID, it returns an ErrNoRecord error.
func (m *UserModel) Get(id int) (User, error) {
    var user User

    stmt := `SELECT id, name, email, created FROM user WHERE id = ?`

    err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return User{}, ErrNoRecord
        } else {
            return User{}, err
        }
    }

    return user, nil
}

func (m *UserModel) UpdatePassword(id int, currentPassword, newPassword string) error {
    var currentHashedPassword []byte

    stmt := `SELECT hashed_password FROM user WHERE id = ?`

    err := m.DB.QueryRow(stmt, id).Scan(&currentHashedPassword)
    if err != nil {
        return err
    }

    err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(currentPassword))
    if err != nil {
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return ErrInvalidCredentials
        } else {
            return err
        }
    }

    newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
    if err != nil {
        return err
    }

    stmt = `UPDATE user SET hashed_password = ? WHERE id = ?`

    _, err = m.DB.Exec(stmt, string(newHashedPassword), id)
    return err
}