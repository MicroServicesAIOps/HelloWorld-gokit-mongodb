package main

import (
	"HelloWorld-gokit-mongodb/api"
	"HelloWorld-gokit-mongodb/db"
	"HelloWorld-gokit-mongodb/db/mongodb"
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
	flag.Parse()

	errc := make(chan error)

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

	// Service domain.
	var service api.Service
	{
		service = api.NewFixedService()
	}

	// Endpoint domain.
	endpoints := api.MakeEndpoints(service)

	// HTTP router
	handler := api.MakeHTTPHandler(endpoints)

	// Create and launch the HTTP server.
	go func() {
		errc <- http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
}
