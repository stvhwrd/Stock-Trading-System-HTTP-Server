package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// ServerConfiguration holds information about a server's network configuration
type ServerConfiguration struct {
	ipAddress string
	port      string
}

// ServerConfigurations holds information about the system's servers' network configurations
type ServerConfigurations struct {
	database    ServerConfiguration
	logging     ServerConfiguration
	transaction ServerConfiguration
	web         ServerConfiguration
}

var serverConfig = ServerConfigurations{}

func init() {
	// Parse and process environment variables from config file
	configurationFilename := ".env"

	err := godotenv.Load(configurationFilename)
	if err != nil {
		log.Fatalf("Error loading %q config file", configurationFilename)
	}

	serverConfig.database.ipAddress = os.Getenv("DATABASE_IP_ADDRESS")
	serverConfig.database.port = os.Getenv("DATABASE_PORT")
	serverConfig.logging.ipAddress = os.Getenv("LOGGING_IP_ADDRESS")
	serverConfig.logging.port = os.Getenv("LOGGING_PORT")
	serverConfig.transaction.ipAddress = os.Getenv("TRANSACTION_IP_ADDRESS")
	serverConfig.transaction.port = os.Getenv("TRANSACTION_PORT")
	serverConfig.web.ipAddress = os.Getenv("WEB_IP_ADDRESS")
	serverConfig.web.port = os.Getenv("WEB_PORT")
}

func main() {
	log.Printf("Server starting with config: \n%+v\n\n", serverConfig)

	// Fire up server
	log.Printf("HTTP server listening on http://%s:%s/\n\n", serverConfig.web.ipAddress, serverConfig.web.port)
	http.HandleFunc("/", requestRouter)
	go log.Fatal(http.ListenAndServe(":"+serverConfig.web.port, nil))
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
