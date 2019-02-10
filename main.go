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
	sendLog(buildLog("HTTP server listening on http://localhost:"+string(state.webServerPort),
		commonlib.DebugType,
		commonlib.LogCommandParameter{
			TransactionNum: "0001",
			Command:        "QUOTE"}))

	// TODO: figure out why commonlib.StartServer doesn't work for this
	http.HandleFunc("/", requestRouter)
	go log.Fatal(http.ListenAndServe(":"+portString, nil))
}

// requestRouter routes the request to the appropriate handler based on its HTTP method
func requestRouter(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request\n", r.Method)
	sendLog(buildLog("Received %s request: "+r.Method,
		commonlib.DebugType,
		commonlib.LogCommandParameter{
			TransactionNum: "0001",
			Command:        "QUOTE"}))

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
		sendLog(buildLog(" HTTP method not supported: "+r.Method,
			commonlib.DebugType,
			commonlib.LogCommandParameter{
				TransactionNum: "0001",
				Command:        "QUOTE"}))
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

// commandHandler decodes a JSON command and forwards it appropriately
func commandHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling JSON body of %s request", r.Method)

	var requestBodyJSON = JSONPayload{}

	// Read request body
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		sendLog(buildLog("Error reading request body",
			commonlib.DebugType,
			commonlib.LogCommandParameter{
				TransactionNum: "0001",
				Command:        "QUOTE"}))

		panic(err)
	}
	defer r.Body.Close()

	// Unmarshal JSON directly into JSONPayload struct
	err = json.Unmarshal(requestBody, &requestBodyJSON)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(""))
		sendLog(buildLog("Error unmarshaling request body",
			commonlib.DebugType,
			commonlib.LogCommandParameter{
				TransactionNum: "0001",
				Command:        "QUOTE"}))
		panic(err)
	}

	sendLog(buildLog("Message received was: "+string(requestBody),
		commonlib.DebugType,
		commonlib.LogCommandParameter{
			TransactionNum: "0001",
			Command:        "QUOTE"}))

	// TODO: Use commonlib full implementation of building message
	commandID := getCommandID(requestBodyJSON.Command, requestBodyJSON.UserID)
	parameters := commonlib.CommandParameter{}

	parameters.UserID = requestBodyJSON.UserID
	parameters.Amount = requestBodyJSON.Amount
	parameters.Filename = requestBodyJSON.Filename
	parameters.StockSymbol = requestBodyJSON.StockSymbol

	destinationServer := getDestinationServer(commandID)

	sendLog(buildLog("Forwarding "+requestBodyJSON.Command+" command to "+destinationServer+" with parameters: "+fmt.Sprintf("%+v", parameters),
		commonlib.DebugType,
		commonlib.LogCommandParameter{
			TransactionNum: "0001",
			Command:        "QUOTE"}))

	response, err := commonlib.SendCommand("POST", "application/json", destinationServer, commonlib.GetSendableCommand(commandID, parameters))
	if err != nil {
		sendLog(buildLog("Error sending command",
			commonlib.DebugType,
			commonlib.LogCommandParameter{
				TransactionNum: "0001",
				Command:        "QUOTE"}))
		panic(err)
	}
	log.Printf("Received response:\n%s\n", response)
}

// userInterfaceHandler serves the user interface HTML file
func userInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving user interface")
	http.ServeFile(w, r, "www/index.html")
}
