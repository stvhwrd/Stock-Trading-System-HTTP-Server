package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	commonlib "github.com/kurtd5105/SENG-468-Common-Lib"
	pql "github.com/lib/pq"
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
	log.Printf("LOG: Received %s request to %q endpoint\n", r.Method, r.URL.Path)

	switch r.Method {
	case "GET":
		log.Printf("LOG: GET request to %q not allowed\n", r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
	case "POST":

		// Read request body
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO: handle retries
			log.Println("LOG: Error reading request body")
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
			err := forwardMessageToTransactionServer(commandID, parameters)
			if err != nil {
				log.Println("LOG: Error while forwarding request to Transaction Server")
				panic(err)
			}
		} else {
			log.Println("LOG: Invalid request -- ignore and discard")
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// uiHandler presents an HTML webpage for a user to manually enter commands, manage their account, etc.
func uiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("LOG: Received %s request to %q endpoint\n", r.Method, r.URL.Path)

	switch r.Method {
	// Incoming data from UI
	case "POST":
		log.Println("LOG: Routing POST request to commandHandler")
		commandHandler(w, r)
	// UI requesting HTML
	case "GET":
		html := "www/index.html"
		log.Printf("LOG: Serving %q\n", html)
		http.ServeFile(w, r, html)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// forwardMessageToTransactionServer forwards the received message to the Transaction Server and returns the error
func forwardMessageToTransactionServer(commandID uint8, parameters commonlib.CommandParameter) error {
	sendableCommand := commonlib.GetSendableCommand(commandID, parameters)
	log.Println("LOG: Sent command to Transaction Server: ", sendableCommand)

	response, err := commonlib.SendCommand("GET", "text/plain", serverState.txPort, sendableCommand)
	log.Println("LOG: Received response from Transaction server: ", response)
	return err
}

// HACK: copied functions from other files to duck "undefined" error

// Begin log_export.go

// getDumplogForUser retrieves the transaction history for a specific user from the database and saves it to a logfile
func getDumplogForUser(userid string) {
	// TODO: Open connection to DB
	dbConnection := pql.ListenerEventConnected
	// TODO: Query DB eg. db.Query("SELECT * FROM transactions WHERE userid = $1", userid)
	// TODO: Write returned rows to <filename>
	// TODO: Close connection to DB
}

// getDumplogForAll retrieves the transaction history for all users from the database and saves it to a logfile
func getDumplogForAll() {
	// TODO: Open connection to DB
	// TODO: Query DB eg. db.Query("SELECT * FROM transactions)
	// TODO: Write returned rows to <filename>
	// TODO: Close connection to DB
}

// saveToFile exports a byte array to a file on local disk
func saveToFile(content []byte, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytesWritten, err := file.Write(content)
	if err != nil {
		panic(err)
	}
	log.Printf("LOG: Wrote %d bytes to logfile: %q\n", bytesWritten, filename)
}

// End log_export.go

// Begin request_validation.go

// validateCommand ensures that the command is one of the known/valid commands and has its necessary parameters
func validateCommandAndParameters(commandID uint8, parameters commonlib.CommandParameter) bool {
	// TODO: Ensure all parameters pertaining to the specific command are present and valid

	return true
}

// validateAmount ensures that the amount specified in command is valid
func validateAmount(amount string) bool {
	// TODO: Must be non-negative
	// TODO: Must not contain non-numerical characters (including "$")
	// TODO: Must contain two decimal places <---(Do we want to round, or reject, if <> 2 decimal places?)

	return true
}

// validateUserID ensures that the user specified in command is valid
func validateUserID(userID string) bool {
	// TODO: Create user if not exists <--(or does the client issue a "CREATE" command?)

	return true
}

// validateStockSymbol ensures that the stock symbol specified in command is valid
func validateStockSymbol(stockSymbol string) bool {
	// TODO: Must be 1 - 3 alphanumeric, case insensitive

	return true
}

// End request_validation.go
