package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	host = flag.String("h", "", "Host name")
	port = flag.String("p", "8080", "port number")
)

func main() {
	flag.Parse()
	r := mux.NewRouter()
	// non REST routes
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("true"))
	}).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	withLogger := handlers.LoggingHandler(os.Stdout, r)
	address := fmt.Sprintf("%s:%s", *host, *port)
	log.Println("Running on " + address)
	http.ListenAndServe(address, withLogger)
}
