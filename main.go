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
	log.Println("Running on " + address)
	http.ListenAndServe(address, r)
}
