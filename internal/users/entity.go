package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"
)

func New(name, login, password string) (*User, error) {
	now := time.Now()
	u := User{Name: name, Login: login, CreatedAt: now, ModifiedAt: now}
	err := u.SetPassword(password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

type User struct {
	ID         int64
	Name       string
	Login      string
	Password   string
	CreatedAt  time.Time
	ModifiedAt time.Time
	Deleted    bool
	LastLogin  time.Time
}

func (u *User) SetPassword(password string) error {
	if password == "" {
		return errors.New("password is required and can't be blank")
	}

	if len(password) < 6 {
		return errors.New("password must have at least 6 characters")
	}

	u.Password = fmt.Sprintf("%x", md5.Sum([]byte(password)))

	return nil
}
