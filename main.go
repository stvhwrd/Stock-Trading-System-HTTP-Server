package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	commonlib "github.com/kurtd5105/SENG-468-Common-Lib"
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

	// Enforce required flags
	if serverState.dbPort == -1 || serverState.txPort == -1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	// Display server state values / config
	log.Printf("LOG: Starting server with config: %+v\n", serverState)

	// Using multiplexer to map endpoints to handler functions
	mux := http.NewServeMux()
	mux.HandleFunc("/ui", uiHandler)
	mux.HandleFunc("/", commandHandler)

	log.Printf("LOG: HTTP server listening on http://localhost:%d\n", serverState.httpPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(serverState.httpPort), mux))
}

// commandHandler processes a JSON command and forwards it to the transaction server
func commandHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("LOG: Received %s request to \"/\" endpoint\n", r.Method)

	switch r.Method {
	case "POST":
		// Read request body
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}
		defer r.Body.Close()

		// Parse and validate request
		commandID, parameters := commonlib.GetCommandFromMessage(requestBody)
		validRequest := validateCommandAndParameters(commandID, parameters)

		if validRequest {
			log.Println("LOG: Forwarding request to Transaction Server")
			// TODO: handle retries
			_, err := forwardMessageToTransactionServer(commandID, parameters)
			if err != nil {
				panic(err)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		w.WriteHeader(http.StatusBadGateway)
	}
}

// uiHandler presents an HTML webpage for a user to manually enter commands, manage their account, etc.
func uiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("LOG: Received %s request to \"/ui\" endpoint\n", r.Method)

	switch r.Method {
	// Incoming data from UI
	case "POST":
		commandHandler(w, r)
		// UI requesting HTML
	case "GET":
		http.ServeFile(w, r, "www/index.html")
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

// forwardMessageToTransactionServer forwards the received message to the Transaction Server
func forwardMessageToTransactionServer(commandID uint8, parameters commonlib.CommandParameter) (bool, error) {
	sendableCommand := commonlib.GetSendableCommand(commandID, parameters)
	log.Println("LOG: Sent command: ", sendableCommand)

	response, err := commonlib.SendCommand("GET", "application/json", serverState.txPort, sendableCommand)
	log.Println("LOG: Received response: ", response)
	if err != nil {
		log.Println("LOG: Forwarding message to Transaction Server failed with error: ", err)
		return false, err
	}

	log.Println("LOG: Forwarding message to Transaction Server succeeded.")
	return true, nil

}
