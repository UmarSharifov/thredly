package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

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
