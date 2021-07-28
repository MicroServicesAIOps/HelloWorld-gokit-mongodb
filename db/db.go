package db

import (
	"HelloWorld-gokit-mongodb/users"
	"errors"
	"flag"
	"fmt"
	"os"
)

type Database interface {
	Init() error
	GetUser(string) (users.User, error)
	GetUsers() ([]users.User, error)
	CreateUser(*users.User) error
	Delete(string, string) error
	Ping() error
}

var (
	database              string
	DefaultDb             Database
	DBTypes               = map[string]Database{}
	ErrNoDatabaseFound    = "No database with name %v registered"
	ErrNoDatabaseSelected = errors.New("No DB selected")
)

func init() {
	flag.StringVar(&database, "database", os.Getenv("USER_DATABASE"), "Database to use, Mongodb or ...")
}

func Init() error {
	//database = "mongodb"
	if database == "" {
		return ErrNoDatabaseSelected
	}
	err := Set()
	if err != nil {
		return err
	}
	return DefaultDb.Init()
}
func Set() error {
	if v, ok := DBTypes[database]; ok {
		DefaultDb = v
		return nil
	}
	return fmt.Errorf(ErrNoDatabaseFound, database)
}
func Register(name string, db Database) {
	DBTypes[name] = db
}

func CreateUser(u *users.User) error {
	return DefaultDb.CreateUser(u)
}

func GetUsers() ([]users.User, error) {
	us, err := DefaultDb.GetUsers()

	return us, err
}
