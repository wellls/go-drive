package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"
)

var (
	BlankPassword       = fmt.Sprintf("%x", md5.Sum([]byte("")))
	ErrNameRequired     = errors.New("name is required")
	ErrLoginRequired    = errors.New("login is required")
	ErrPasswordRequired = errors.New("password is required and can't be blank")
	ErrPasswordLen      = errors.New("password must have at least 6 characters")
)

func New(name, login, password string) (*User, error) {
	u := User{Name: name, Login: login, ModifiedAt: time.Now()}
	err := u.SetPassword(password)
	if err != nil {
		return nil, err
	}

	err = u.Validate()
	if err != nil {
		return nil, err
	}

	return &u, nil
}

type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
	LastLogin  time.Time `json:"last_login"`
}

func (u *User) SetPassword(password string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < 6 {
		return ErrPasswordLen
	}

	u.Password = fmt.Sprintf("%x", md5.Sum([]byte(password)))

	return nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameRequired
	}

	if u.Login == "" {
		return ErrLoginRequired
	}

	if u.Password == BlankPassword {
		return ErrPasswordRequired
	}

	return nil
}
