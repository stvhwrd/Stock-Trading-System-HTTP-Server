package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// ServerNetwork holds information about the system's servers' network addresses
type ServerNetwork struct {
	databaseServerAddressAndPort    string
	loggingServerAddressAndPort     string
	transactionServerAddressAndPort string
	webServerAddress                string
	webServerPort                   string
}

var serverNetworkConfig = ServerNetwork{}

func init() {
	// Parse and process CLI flags
	flag.StringVar(&serverNetworkConfig.databaseServerAddressAndPort, "db", "", "[REQUIRED] the IP address and port on which the USER ACCOUNT DATABASE server is running, eg. -db=localhost:8080")
	flag.StringVar(&serverNetworkConfig.loggingServerAddressAndPort, "log", "", "[REQUIRED] the IP address and port on which the LOGGING DATABASE server is running, eg. -log=localhost:8081")
	flag.StringVar(&serverNetworkConfig.webServerAddress, "webip", "", "[REQUIRED] the IP address on which *this* HTTP server is running, eg. -webip=localhost")
	flag.StringVar(&serverNetworkConfig.webServerPort, "webport", "", "[REQUIRED] the IP address and port on which *this* HTTP server is running, eg. -webport=localhost:80")
	flag.StringVar(&serverNetworkConfig.transactionServerAddressAndPort, "tx", "", "[REQUIRED] the IP address and port on which the TRANSACTION server is running, eg. -tx=localhost:8082")
	flag.Parse()

	// Enforce required flags
	if serverNetworkConfig.databaseServerAddressAndPort == "" ||
		serverNetworkConfig.transactionServerAddressAndPort == "" ||
		serverNetworkConfig.loggingServerAddressAndPort == "" ||
		serverNetworkConfig.webServerAddress == "" ||
		serverNetworkConfig.webServerPort == "" {
		log.Println("Error: Required flags were not provided at runtime")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	log.Printf("Server starting with config: \n%+v\n\n", serverNetworkConfig)

	// Fire up server
	log.Printf("HTTP server listening on http://%s/\n\n", serverNetworkConfig.webServerAddress)
	// commonlib.StartServer(serverNetworkConfig.webServerPort, requestRouter)
	http.HandleFunc("/", requestRouter)
	go log.Fatal(http.ListenAndServe(":"+serverNetworkConfig.webServerPort, nil))
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

// requestDecoder decodes a JSON command
func commandHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling JSON body of %s request", r.Method)

	var requestBodyJSON = JSONPayload{}

	// Read request body
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panic(err)
	}
	defer r.Body.Close()

	// Unmarshal JSON directly into JSONPayload struct
	err = json.Unmarshal(requestBody, &requestBodyJSON)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panic(err)
	}

	log.Printf("Message received was: %+v", requestBodyJSON)
	// TODO: Validate parameters
	buildAndSendMessage(requestBodyJSON)
}

// userInterfaceHandler serves the user interface HTML file
func userInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving user interface")
	http.ServeFile(w, r, "www/index.html")
}
