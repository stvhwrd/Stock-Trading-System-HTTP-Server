package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// serverConfig holds information about the server's state and configuration
type serverState struct {
	databasePort          int
	loggingPort           int
	transactionServerPort int
	httpServerPort        int
}

var serverConfig = serverState{}

func init() {
	// Parse and process CLI flags (type is implicit)
	flag.IntVar(&serverConfig.databasePort, "dbport", -1, "[REQUIRED] the port on which the USER ACCOUNT DATABASE server is running, eg. -dbport=8080")
	flag.IntVar(&serverConfig.loggingPort, "logport", -1, "[REQUIRED] the port on which the LOGGING DATABASE server is running, eg. -logport=8081")
	flag.IntVar(&serverConfig.httpServerPort, "httpport", 80, "[optional] the port on which *this* HTTP server is running, eg. -httpport=80")
	flag.IntVar(&serverConfig.transactionServerPort, "txport", -1, "[REQUIRED] the port on which the TRANSACTION server is running, eg. -txport=8082")
	flag.Parse()

	// Enforce required flags
	if serverConfig.databasePort == -1 || serverConfig.transactionServerPort == -1 || serverConfig.loggingPort == -1 {
		flag.PrintDefaults()
		log.Fatal("Flags not provided at runtime")
	}
}

func main() {
	log.Printf("Server starting with config: %+v\n", serverConfig)
	// router.HandleFunc("/create", CreateEndpoint).Methods("GET")
	http.HandleFunc("/", requestRouter)

	// Fire up server
	log.Printf("HTTP server listening on http://localhost:%d/\n", serverConfig.httpServerPort)
	go log.Fatal(http.ListenAndServe(":"+strconv.Itoa(serverConfig.httpServerPort), nil))
}

// requestRouter routes the request to the appropriate handler based on its HTTP method
func requestRouter(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request\n", r.Method)

	switch r.Method {
	case http.MethodPost:
		log.Println("Routing POST request to commandHandler")
		commandHandler(w, r)
		// GET requests are only expected from UI
	case http.MethodGet:
		log.Println("Routing GET request to userInterfaceHandler")
		userInterfaceHandler(w, r)
	default:
		// No other HTTP methods are supported
		log.Printf("HTTP method %q not supported\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ExpectedJSON represents the expected JSON body of a request
type ExpectedJSON struct {
	Command     string `json: "command"`
	UserID      string `json: "userID"`
	Amount      string `json: "amount"`
	StockSymbol string `json: "stockSymbol"`
	Filename    string `json: "filename"`
}

var requestBodyJSON = ExpectedJSON{}

// commandHandler processes a JSON command and forwards it to the transaction server
func commandHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling JSON body of %s request", r.Method)

	requestBody, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panic(readError)
	}
	defer r.Body.Close()

	unmarshalError := json.Unmarshal(requestBody, &requestBodyJSON)
	if unmarshalError != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Panic(unmarshalError)
	}

	log.Printf("Message received was: %+v", requestBodyJSON)
	// Parse and validate request
	// Call appropriate function based
}

// userInterfaceHandler serves the user interface HTML file
func userInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving user interface")
	http.ServeFile(w, r, "www/index.html")
}
