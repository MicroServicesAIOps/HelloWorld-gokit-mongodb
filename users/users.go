package users

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"time"
)

type User struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"-" bson:"email"`
	Username  string `json:"username" bson:"username"`
	Password  string `json:"-" bson:"password,omitempty"`
	UserID    string `json:"id" bson:"-"`
	Salt      string `json:"-" bson:"salt"`
}

func New() User {
	u := User{}
	u.NewSalt()
	return u
}
func (u *User) NewSalt() {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	u.Salt = fmt.Sprintf("%x", h.Sum(nil))
}
