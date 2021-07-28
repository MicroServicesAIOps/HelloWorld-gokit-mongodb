package api

import (
	"HelloWorld-gokit-mongodb/db"
	"HelloWorld-gokit-mongodb/users"
	"crypto/sha1"
	"fmt"
	"io"
	"time"
)

type IMyService interface {
	Register(username, password, email, first, last string) (string, error)
	GetUsers(id string) ([]users.User, error)
	Delete(entity, id string) error
	Health() []Health // GET /health
}

type MyService struct{}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func (s MyService) Register(username, password, email, first, last string) (string, error) {
	u := users.New()
	u.Username = username
	u.Password = calculatePassHash(password, u.Salt)
	u.Email = email
	u.FirstName = first
	u.LastName = last
	err := db.CreateUser(&u)
	return u.UserID, err
}

func (s MyService) GetUsers(id string) ([]users.User, error) {
	if id == "" {
		us, err := db.GetUsers()
		return us, err
	}
	u, err := db.GetUser(id)
	return []users.User{u}, err
}

func (s MyService) Delete(entity, id string) error {
	return db.Delete(entity, id)
}

func (s MyService) Health() []Health {
	var health []Health
	app := Health{"user", "OK", time.Now().String()}

	health = append(health, app)

	return health
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}
