package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	host = flag.String("h", "", "Host name")
	port = flag.String("p", "8080", "port number")
)

func main() {
	flag.Parse()
	address := fmt.Sprintf("%s:%s", *host, *port)
	r := mux.NewRouter()
	// non REST routes
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))
	r.PathPrefix("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("true"))
	}).Methods("GET")

	log.Println("Running on " + address)
	http.ListenAndServe(address, r)
}
