package main

import (
	"HelloWorld-gokit-mongodb/api"
	"HelloWorld-gokit-mongodb/db"
	"HelloWorld-gokit-mongodb/db/mongodb"
	"flag"
	corelog "log"
	"net/http"

	httpTransport "github.com/go-kit/kit/transport/http"
)

var (
	port string
)

const (
	ServiceName = "user"
)

func init() {
	flag.StringVar(&port, "port", "8084", "Port on which to run")
	db.Register("mongodb", &mongodb.Mongo{})
}

func main() {

	dbconn := false
	for !dbconn {
		err := db.Init()
		if err != nil {
			if err == db.ErrNoDatabaseSelected {
				corelog.Fatal(err)
			}
			corelog.Print(err)
		} else {
			dbconn = true
		}
	}

	s := api.MyService{}
	health := api.MakeHealthEndpoint(s)
	regist := api.MakeRegisterEndpoint(s)
	userGet := api.MakeUserGetEndpoint(s)
	delete := api.MakeDeleteEndpoint(s)

	healthServer := httpTransport.NewServer(health, api.HealthDecodeRequest, api.HealthEncodeResponse)

	registServer := httpTransport.NewServer(regist, api.DecodeRegisterRequest, api.EncodeResponse)
	userGetServer := httpTransport.NewServer(userGet, api.DecodeGetRequest, api.EncodeResponse)
	deleteServer := httpTransport.NewServer(delete, api.DecodeDeleteRequest, api.EncodeResponse)

	http.Handle("/regist/", registServer)
	http.Handle("/userget/", userGetServer)
	http.Handle("/delete/", deleteServer)
	http.Handle("/health", healthServer)
	go http.ListenAndServe("0.0.0.0:8084", nil)

	select {}
}
