package api

import (
	"HelloWorld-gokit-mongodb/db"
	"HelloWorld-gokit-mongodb/users"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

type Service interface {
	Register(username, password, email, first, last string) (string, error)
	GetUsers(id string) ([]users.User, error)
	Delete(entity, id string) error
	Health() []Health // GET /health
}

func NewFixedService() Service {
	return &fixedService{}
}

type fixedService struct{}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func (s *fixedService) Register(username, password, email, first, last string) (string, error) {
	u := users.New()
	u.Username = username
	u.Password = calculatePassHash(password, u.Salt)
	u.Email = email
	u.FirstName = first
	u.LastName = last
	err := db.CreateUser(&u)
	return u.UserID, err
}

func (s *fixedService) GetUsers(id string) ([]users.User, error) {
	if id == "" {
		us, err := db.GetUsers()
		return us, err
	}
	u, err := db.GetUser(id)
	return []users.User{u}, err
}

func (s *fixedService) Delete(entity, id string) error {
	return db.Delete(entity, id)
}

func (s *fixedService) Health() []Health {
	var health []Health
	dbstatus := "OK"

	err := db.Ping()
	if err != nil {
		dbstatus = "err"
	}

	app := Health{"user", "OK", time.Now().String()}
	db := Health{"user-db", dbstatus, time.Now().String()}

	health = append(health, app)
	health = append(health, db)

	return health
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}