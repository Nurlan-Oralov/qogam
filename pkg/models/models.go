package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials ошибка. Мы будем использовать это позже, если пользователь
	// пытается войти в систему с неправильным адресом электронной почты или паролем.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail ошибка. Мы будем использовать это позже, если пользователь
	// пытается зарегистрироваться с помощью адреса электронной почты, который уже используется.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
