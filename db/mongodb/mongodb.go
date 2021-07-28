package mongodb

import (
	"HelloWorld-gokit-mongodb/users"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	name     string
	password string
	host     string
	db       = "users"
	//ErrInvalidHexID represents a entity id that is not a valid bson ObjectID
	ErrInvalidHexID = errors.New("Invalid Id Hex")
)

func init() {
	flag.StringVar(&name, "mongo-user", os.Getenv("MONGO_USER"), "Mongo user")
	flag.StringVar(&password, "mongo-password", os.Getenv("MONGO_PASS"), "Mongo password")
	flag.StringVar(&host, "mongo-host", os.Getenv("MONGO_HOST"), "Mongo host")
	fmt.Println(name)
	fmt.Println(password)
	fmt.Println(host)
}

// Mongo meets the Database interface requirements
type Mongo struct {
	//Session is a MongoDB Session
	Session *mgo.Session
}

// Init MongoDB
func (m *Mongo) Init() error {
	u := getURL()
	var err error
	m.Session, err = mgo.DialWithTimeout(u.String(), time.Duration(5)*time.Second)
	if err != nil {
		return err
	}
	return m.EnsureIndexes()
}

// MongoUser is a wrapper for the users
type MongoUser struct {
	users.User `bson:",inline"`
	ID         bson.ObjectId `bson:"_id"`
}

// New Returns a new MongoUser
func New() MongoUser {
	u := users.New()
	return MongoUser{
		User: u,
	}
}

// AddUserIDs adds userID as string to user
func (mu *MongoUser) AddUserIDs() {
	mu.User.UserID = mu.ID.Hex()
}

// CreateUser Insert user to MongoDB, including connected addresses and cards, update passed in user with Ids
func (m *Mongo) CreateUser(u *users.User) error {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := New()
	mu.User = *u
	mu.ID = id

	c := s.DB("").C("customers")
	_, err := c.UpsertId(mu.ID, mu)
	if err != nil {
		return err
	}
	mu.User.UserID = mu.ID.Hex()

	*u = mu.User
	return nil
}

// GetUser Get user by their object id
func (m *Mongo) GetUser(id string) (users.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	if !bson.IsObjectIdHex(id) {
		return users.New(), errors.New("Invalid Id Hex")
	}
	c := s.DB("").C("customers")
	mu := New()
	err := c.FindId(bson.ObjectIdHex(id)).One(&mu)
	mu.AddUserIDs()
	return mu.User, err
}

// GetUsers Get all users
func (m *Mongo) GetUsers() ([]users.User, error) {
	// TODO: add paginations
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("customers")
	var mus []MongoUser
	err := c.Find(nil).All(&mus)
	us := make([]users.User, 0)
	for _, mu := range mus {
		mu.AddUserIDs()
		us = append(us, mu.User)
	}
	return us, err
}

// CreateAddress Inserts Address into MongoDB
func (m *Mongo) Delete(entity, id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Invalid Id Hex")
	}
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C(entity)
	return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}

func getURL() url.URL {
	ur := url.URL{
		Scheme: "mongodb",
		Host:   host,
		Path:   db,
	}
	if name != "" {
		u := url.UserPassword(name, password)
		ur.User = u
	}
	return ur
}

// EnsureIndexes ensures username is unique
func (m *Mongo) EnsureIndexes() error {
	s := m.Session.Copy()
	defer s.Close()
	i := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	c := s.DB("").C("customers")
	return c.EnsureIndex(i)
}

func (m *Mongo) Ping() error {
	s := m.Session.Copy()
	defer s.Close()
	return s.Ping()
}
