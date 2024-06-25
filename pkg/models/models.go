package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")
var ErrInvalidCredentials = errors.New("models: неверные учетные данные")

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
	ULastName       string
	UFirstName      string
	UPhoto          string
	PublicationDate time.Time
	ViewsCount      int
	Content         string
	Photo           string
	ParentId        sql.NullInt64
	Category        string
	ChildTreds      []*Tred
	Tags            []string
}

type TredsWithChilds struct {
	MainTred   Tred
	ChildTreds []*Tred
}

type Categories struct {
	ID   int
	Name string
}

type Events struct {
	ID              int
	UserId          int
	PublicationDate time.Time
	ViewsCount      int
	Content         string
	CategoryId      sql.NullInt64
	Photo           string
}

type EventCategories struct {
	ID   int
	Name string
}

type Complaint struct {
	ID            int
	UserId        int
	TredId        int
	ComplaintDate time.Time
	Description   string
}

type ComplaintWithDetails struct {
	ID            int
	TredID        int
	UserID        int
	ComplaintDate time.Time
	UserFirstName string
	UserLastName  string
	TredContent   string
}

type Subscribe struct {
	id                int
	subscriber_id     int
	subscribed_to_id  int
	subscription_date time.Time
}

type Tags struct {
	id        int
	thread_id int
	tag_name  string
}
