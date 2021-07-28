package main

import (
	"HelloWorld-gokit-mongodb/api"
	"HelloWorld-gokit-mongodb/db"
	"HelloWorld-gokit-mongodb/db/mongodb"
	"context"
	"flag"
	"fmt"
	corelog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
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
	ctx := context.Background()
	errChan := make(chan error)

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

	endpoints := api.Endpoints{
		HealthEndpoint:   health,
		RegisterEndpoint: regist,
		UserGetEndpoint:  userGet,
		DeleteEndpoint:   delete,
	}
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	r := api.MakeHttpHandler(ctx, endpoints, logger)

	go func() {
		fmt.Println("Http Server start at port:8084")
		handler := r
		errChan <- http.ListenAndServe(":8084", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(<-errChan)

	// healthServer := httpTransport.NewServer(health, api.HealthDecodeRequest, api.HealthEncodeResponse)

	// registServer := httpTransport.NewServer(regist, api.DecodeRegisterRequest, api.EncodeResponse)
	// userGetServer := httpTransport.NewServer(userGet, api.DecodeUserRequest, api.EncodeResponse)
	// deleteServer := httpTransport.NewServer(delete, api.DecodeDeleteRequest, api.EncodeResponse)

	// http.Handle("/regist/", registServer)
	// http.Handle("/userget/", userGetServer)
	// http.Handle("/delete/", deleteServer)
	// http.Handle("/health", healthServer)
	// go http.ListenAndServe("0.0.0.0:8084", nil)

	// select {}
}
