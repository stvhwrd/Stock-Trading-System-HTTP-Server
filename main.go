package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

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
	flag.StringVar(&state.databaseServerAddressAndPort, "db", "",
		"[REQUIRED] the IP address and port on which the USER ACCOUNT DATABASE server is running, eg. -db=localhost:8080")
	flag.StringVar(&state.loggingServerAddressAndPort, "log", "",
		"[REQUIRED] the IP address and port on which the LOGGING DATABASE server is running, eg. -log=localhost:8081")
	flag.IntVar(&state.webServerPort, "port", -1,
		"[REQUIRED] the port on which *this* HTTP server is running, eg. -port=localhost:80")
	flag.StringVar(&state.transactionServerAddressAndPort, "tx", "",
		"[REQUIRED] the IP address and port on which the TRANSACTION server is running, eg. -tx=localhost:8082")
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

	commonlib.ServerName = "http-server"
}

func main() {
	portString := strconv.Itoa(state.webServerPort)

	// Fire up server
	log.Printf("HTTP server listening on http://localhost:%d/\n\n", state.webServerPort)

	// TODO: figure out why commonlib.StartServer doesn't work for this
	http.HandleFunc("/", requestRouter)
	log.Fatal(http.ListenAndServe(":"+portString, nil))
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
		// No other HTTP methods are supported
		errorMessage := fmt.Sprintf("HTTP method not supported: %s\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(errorMessage))
		log.Fatal(errorMessage)
	}
}

// JSONPayload represents the expected JSON body of a request
type JSONPayload struct {
	Message     string `json: "message"` // HACK
	UserID      string `json: "userID,omitempty"`
	Amount      string `json: "amount,omitempty"`
	StockSymbol string `json: "stockSymbol,omitempty"`
	Filename    string `json: "filename,omitempty"`
}

// commandHandler decodes a JSON command and forwards it appropriately
func commandHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling JSON body of %s request", r.Method)

	// Read request body
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("Error reading request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		log.Fatalln(errorMessage)
	}
	defer r.Body.Close()

	// Unmarshal JSON directly into JSONPayload struct
	var requestBodyJSON = JSONPayload{}

	if err = json.Unmarshal(requestBody, &requestBodyJSON); err != nil {
		errorMessage := fmt.Sprintf("Error unmarshaling request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		log.Fatalln(errorMessage)
	}

	//Increment global transaction counter
	transactionNumString := strconv.FormatUint(incrementTransactionNum(), 10)

	// Extract commandID from message
	message, err := strconv.ParseInt(requestBodyJSON.Message, 10, 8)
	if err != nil {
		errorMessage := fmt.Sprintf("Error parsing message content: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		log.Fatalln(errorMessage)
	}
	commandID := uint8(message)

	// Build a CommandParameter to send to Transaction Server
	parameters := commonlib.CommandParameter{
		UserID:         requestBodyJSON.UserID,
		Amount:         requestBodyJSON.Amount,
		Filename:       requestBodyJSON.Filename,
		StockSymbol:    requestBodyJSON.StockSymbol,
		TransactionNum: transactionNumString,
	}

	// Build a LogCommandParameter to send to Logging Server
	loggingParameters := commonlib.LogCommandParameter{
		Username:       requestBodyJSON.UserID,
		Funds:          requestBodyJSON.Amount,
		LogFilename:    requestBodyJSON.Filename,
		LogStockSymbol: requestBodyJSON.StockSymbol,
		Server:         "Web",
		TransactionNum: transactionNumString,
		Timestamp:      commonlib.GetTimeStampString(),
		Command:        commonlib.CommandNames[commandID],
	}

	sendLog(buildLog(fmt.Sprintf("Received request: %s", requestBodyJSON),
		commonlib.DebugType,
		loggingParameters))

	// Destination depends on type of command
	destinationServer := getDestinationServer(commandID)

	sendLog(buildLog(fmt.Sprintf("Forwarding #%s command to %s with parameters: %+v\n",
		requestBodyJSON.Message, destinationServer, parameters),
		commonlib.SystemEventType,
		loggingParameters))

	response, err := commonlib.SendCommand(
		"POST",
		"application/json",
		destinationServer,
		commonlib.GetSendableCommand(commandID, parameters))

	if err != nil {
		errorMessage := fmt.Sprintf("Error sending command: %s\n\n Server response: %s\n",
			err.Error(), response)
		sendLog(buildLog(
			errorMessage,
			commonlib.ErrorEventType,
			loggingParameters))
		log.Fatalf(errorMessage)
	}

	sendLog(buildLog(
		fmt.Sprintf("%s responded: %s\n", destinationServer, response),
		commonlib.DebugType,
		loggingParameters))

	// Request received intact
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// userInterfaceHandler serves the user interface HTML file
func userInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving user interface")
	http.ServeFile(w, r, "www/index.html")
}

// Global transactionNum initialized to 0 by default
var transactionNum uint64

// incrementTransactionNum atomically increments the global transaction counter and returns its value
func incrementTransactionNum() uint64 {
	return atomic.AddUint64(&transactionNum, 1)
}
