package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	commonlib "github.com/kurtd5105/SENG-468-Common-Lib"
)

// ServerNetwork holds information about the system's servers' network addresses
type ServerNetwork struct {
	databaseServerAddressAndPort    string
	loggingServerAddressAndPort     string
	transactionServerAddressAndPort string
	webServerPort                   int
}

var state = ServerNetwork{}

func init() {
	// Parse and process CLI flags
	flag.StringVar(&state.databaseServerAddressAndPort, "db", "", "[REQUIRED] the IP address and port on which the USER ACCOUNT DATABASE server is running, eg. -db=localhost:8080")
	flag.StringVar(&state.loggingServerAddressAndPort, "log", "", "[REQUIRED] the IP address and port on which the LOGGING DATABASE server is running, eg. -log=localhost:8081")
	flag.IntVar(&state.webServerPort, "port", -1, "[REQUIRED] the port on which *this* HTTP server is running, eg. -port=localhost:80")
	flag.StringVar(&state.transactionServerAddressAndPort, "tx", "", "[REQUIRED] the IP address and port on which the TRANSACTION server is running, eg. -tx=localhost:8082")
	flag.Parse()

	// Enforce required flags
	if state.databaseServerAddressAndPort == "" ||
		state.transactionServerAddressAndPort == "" ||
		state.loggingServerAddressAndPort == "" ||
		state.webServerPort < 0 {
		log.Println("Error: Required flags were not provided at runtime")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	log.Printf("Server starting with config: \n%+v\n\n", state)

	portString := strconv.Itoa(state.webServerPort)

	// Fire up server
	log.Printf("HTTP server listening on http://localhost:%d/\n\n", state.webServerPort)

	// TODO: figure out why commonlib.StartServer doesn't work for this
	http.HandleFunc("/", requestRouter)
	go log.Fatal(http.ListenAndServe(":"+portString, nil))
}

// requestRouter routes the request to the appropriate handler based on its HTTP method
func requestRouter(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request\n", r.Method)

	switch r.Method {
	case http.MethodPost:
		// POST requests come from UI and/or workload generator:
		log.Println("Routing POST request to commandHandler")
		commandHandler(w, r)
	// GET requests are only expected from UI:
	case http.MethodGet:
		log.Println("Routing GET request to userInterfaceHandler")
		userInterfaceHandler(w, r)
	default:
		// No other HTTP methods are supported:
		log.Printf("HTTP method %q not supported\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// JSONPayload represents the expected JSON body of a request
type JSONPayload struct {
	Command     string `json: "command"`
	UserID      string `json: "userID,omitempty"`
	Amount      string `json: "amount,omitempty"`
	StockSymbol string `json: "stockSymbol,omitempty"`
	Filename    string `json: "filename,omitempty"`
}

// requestDecoder decodes a JSON command and forwards it appropriately
func commandHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling JSON body of %s request", r.Method)

	var requestBodyJSON = JSONPayload{}

	// Read request body
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading request body"))
		log.Panic(err)
	}
	defer r.Body.Close()

	// Unmarshal JSON directly into JSONPayload struct
	err = json.Unmarshal(requestBody, &requestBodyJSON)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error unmarshaling request body"))
		log.Panic(err)
	}

	log.Printf("Message received was: %+v", requestBodyJSON)

	// TODO: Use commonlib full implementation of building message
	commandID := getCommandID(requestBodyJSON.Command, requestBodyJSON.UserID)
	parameters := commonlib.CommandParameter{}

	parameters.UserID = requestBodyJSON.UserID
	parameters.Amount = requestBodyJSON.Amount
	parameters.Filename = requestBodyJSON.Filename
	parameters.StockSymbol = requestBodyJSON.StockSymbol

	destinationServer := getDestinationServer(commandID)

	log.Printf("Sending POST request with command to %s/ \n\n", destinationServer)
	log.Printf("\tCOMMAND: %s (#%d)\n", requestBodyJSON.Command, commandID)
	log.Printf("\tPARAMETERS: %+v\n", parameters)

	response, err := commonlib.SendCommand("POST", "application/json", destinationServer, commonlib.GetSendableCommand(commandID, parameters))
	if err != nil {
		log.Printf("-- Error sending command --")
		panic(err)
	}
	log.Printf("Received response:\n%s\n", response)
}

// userInterfaceHandler serves the user interface HTML file
func userInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving user interface")
	http.ServeFile(w, r, "www/index.html")
}
