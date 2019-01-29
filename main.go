package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// ServerState holds information about the server's state and configuration
type ServerState struct {
	dbPort   int
	txPort   int
	httpPort int
}

var serverState = ServerState{}

func init() {
	// Parse and process CLI flags
	flag.IntVar(&serverState.dbPort, "dbport", -1, "[REQUIRED] the port on which the DATABASE server is running, eg. --dbport=8080")
	flag.IntVar(&serverState.httpPort, "httpport", 80, "[optional -- default is port 80] the port on which *this* HTTP server is running, eg. --httpport=80")
	flag.IntVar(&serverState.txPort, "txport", -1, "[REQUIRED] the port on which the TRANSACTION server is running, eg. --txport=8082")
	flag.Parse()

	// Force flags as required
	if serverState.dbPort == -1 || serverState.txPort == -1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	// Display server state values
	fmt.Printf("%+v\n", serverState)

	// Fire up the server
	portNumString := ":" + strconv.Itoa(serverState.httpPort)
	http.HandleFunc("/", requestHandler)
	http.HandleFunc("/ui", uiHandler)
	fmt.Printf("HTTP server listening on http://localhost:%d\n", serverState.httpPort)
	http.ListenAndServe(portNumString, nil)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HTTP server running right meow.\n\n")
	fmt.Fprintf(w, "The request received:\n\n")
	fmt.Fprintf(w, "%s\n", r.Method)
	fmt.Fprintf(w, "%s\n", r.Body)
}

// uiHandler presents an HTML webpage for a user to manually enter commands, manage their account, etc.
func uiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("LOG: Received %s request to \"/ui\" endpoint\n", r.Method)

	switch r.Method {
	case "POST":
		requestHandler(w, r)
	case "GET":
		http.ServeFile(w, r, "www/index.html")
	}
}
