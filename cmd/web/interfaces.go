package main

import "snippetbox/internal/models"

type userModelInterface interface {
    Insert(name, email, password string) error
    Get(id int) (models.User, error)
    Exists(id int) (bool, error)
    Authenticate(email, password string) (int, error)
    UpdatePassword(id int, currentPassword, newPassword string) error
}

type snippetModelInterface interface {
    Insert(title string, content string, expires int) (int, error)
    Get(id int) (models.Snippet, error)
    Latest(n int) ([]models.Snippet, error)
}