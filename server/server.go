package main

import (
	"flag"
	myhandler "github.com/davejohnston/graphql-go-tutorial/handlers"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func main() {

	flag.Parse()

	router := mux.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.HandleFunc("/subscriptions", myhandler.WebsocketHandler())
	router.HandleFunc("/graphql", myhandler.GraphQLHandler())

	srv := &http.Server{
		Handler:      handlers.CORS(methodsOk, headersOk, originsOk)(router),
		Addr:         "localhost:4000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	glog.Info("Starting GraphQL Server on ", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		glog.Error("Error starting server ", err)
	}
}
