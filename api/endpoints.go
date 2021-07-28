package api

import (
	"HelloWorld-gokit-mongodb/users"
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	HealthEndpoint   endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
	UserGetEndpoint  endpoint.Endpoint
}

type GetRequest struct {
	ID   string
	Attr string
}

type userResponse struct {
	User users.User `json:"user"`
}

type usersResponse struct {
	Users []users.User `json:"customer"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type statusResponse struct {
	Status bool `json:"status"`
}

type postResponse struct {
	ID string `json:"id"`
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

func MakeRegisterEndpoint(s IMyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		id, err := s.Register(req.Username, req.Password, req.Email)
		return postResponse{ID: id}, err
	}
}

func MakeUserGetEndpoint(s IMyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(GetRequest)
		fmt.Print(req)
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

func MakeHealthEndpoint(s IMyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health := s.Health()
		return healthResponse{Health: health}, nil
	}
}
