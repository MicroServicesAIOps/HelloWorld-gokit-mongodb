package api

import (
	"HelloWorld-gokit-mongodb/users"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrInvalidRequest = errors.New("Invalid request")
)

func EncodeError(c context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/hal+json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       err.Error(),
		"status_code": code,
		"status_text": http.StatusText(code),
	})
}

func DecodeRegisterRequest(c context.Context, r *http.Request) (interface{}, error) {
	reg := registerRequest{}
	err := json.NewDecoder(r.Body).Decode(&reg)
	fmt.Println(r)
	fmt.Println(reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func DecodeDeleteRequest(c context.Context, r *http.Request) (interface{}, error) {
	d := deleteRequest{}
	u := strings.Split(r.URL.Path, "/")
	if len(u) == 3 {
		d.Entity = u[1]
		d.ID = u[2]
		return d, nil
	}
	return d, ErrInvalidRequest
}

func DecodeGetRequest(c context.Context, r *http.Request) (interface{}, error) {
	g := GetRequest{}
	fmt.Println(r.URL.Path)
	u := strings.Split(r.URL.Path, "/")
	if len(u) > 2 {
		g.ID = u[2]
		if len(u) > 3 {
			g.Attr = u[3]
		}
	}
	return g, nil
}

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	u := users.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func HealthDecodeRequest(c context.Context, request *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

func HealthEncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func EncodeResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	// All of our response objects are JSON serializable, so we just do that.
	w.Header().Set("Content-Type", "application/hal+json")
	return json.NewEncoder(w).Encode(response)
}
