package api

import (
	"HelloWorld-gokit-mongodb/users"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	RegisterEndpoint endpoint.Endpoint
	UserGetEndpoint  endpoint.Endpoint
	DeleteEndpoint   endpoint.Endpoint
	HealthEndpoint   endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		RegisterEndpoint: MakeRegisterEndpoint(s),
		HealthEndpoint:   MakeHealthEndpoint(s),
		UserGetEndpoint:  MakeUserGetEndpoint(s),
		DeleteEndpoint:   MakeDeleteEndpoint(s),
	}
}

func MakeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		id, err := s.Register(req.Username, req.Password, req.Email, req.FirstName, req.LastName)
		return postResponse{ID: id}, err
	}
}

func MakeUserGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(GetRequest)
		usrs, err := s.GetUsers(req.ID)

		if req.ID == "" {
			return EmbedStruct{usersResponse{Users: usrs}}, err
		}
		if len(usrs) == 0 {
			return users.User{}, err
		}
		user := usrs[0]
		return user, err
	}
}

// MakeLoginEndpoint returns an endpoint via the given service.
func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(deleteRequest)
		err = s.Delete(req.Entity, req.ID)
		if err == nil {
			return statusResponse{Status: true}, err
		}
		return statusResponse{Status: false}, err
	}
}

// MakeHealthEndpoint returns current health of the given service.
func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health := s.Health()
		return healthResponse{Health: health}, nil
	}
}

type GetRequest struct {
	ID   string
	Attr string
}

type loginRequest struct {
	Username string
	Password string
}

type userResponse struct {
	User users.User `json:"user"`
}

type usersResponse struct {
	Users []users.User `json:"customer"`
}

type registerRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type statusResponse struct {
	Status bool `json:"status"`
}

type postResponse struct {
	ID string `json:"id"`
}

type deleteRequest struct {
	Entity string
	ID     string
}

type healthRequest struct {
	//
}

type healthResponse struct {
	Health []Health `json:"health"`
}

type EmbedStruct struct {
	Embed interface{} `json:"_embedded"`
}