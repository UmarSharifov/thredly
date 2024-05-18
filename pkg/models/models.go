package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")
var ErrInvalidCredentials = errors.New("models: неверные учетные данные")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	LastName       string
	FirstName      string
	Photo          string
	DateOfBirthDay time.Time
	Email          string
	PhoneNumber    string
	UserLogin      string
	UserPwd        string
}

type Tred struct {
	ID              int
	UserId          int
	PublicationDate time.Time
	ViewsCount      int
	Content         string
	Photo           string
}
