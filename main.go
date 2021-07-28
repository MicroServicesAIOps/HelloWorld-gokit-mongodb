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

	endpoints := api.Endpoints{
		HealthEndpoint:   health,
		RegisterEndpoint: regist,
		UserGetEndpoint:  userGet,
	}

	r := api.MakeHttpHandler(ctx, endpoints)

	go func() {
		fmt.Println("Http Server start")
		fmt.Println(port)
		handler := r
		errChan <- http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(<-errChan)
}
