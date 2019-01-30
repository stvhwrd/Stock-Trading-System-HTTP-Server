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
	log.Printf("LOG: %+v\n", serverState)

	// Using multiplexer to map endpoints to handler functions
	mux := http.NewServeMux()
	mux.HandleFunc("/ui", uiHandler)
	mux.HandleFunc("/", commandHandler)

	log.Printf("LOG: HTTP server listening on http://localhost:%d\n", serverState.httpPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(serverState.httpPort), mux))
}

// commandHandler processes a JSON command and forwards it to the transaction server
func commandHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("LOG: Received %s request to \"/ui\" endpoint\n", r.Method)

	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusBadRequest)
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
		if validateCommandAndParameters(commandID, parameters) {
			forwardMessageToTransactionServer(commandID, parameters)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// uiHandler presents an HTML webpage for a user to manually enter commands, manage their account, etc.
func uiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("LOG: Received %s request to \"/ui\" endpoint\n", r.Method)

	switch r.Method {
	case "POST":
		commandHandler(w, r)
	case "GET":
		http.ServeFile(w, r, "www/index.html")
	}
}

// getDumplogForUser retrieves the transaction history for a specific user from the database and saves it to a logfile
func getDumplogForUser(userid string) {
	// TODO: Open connection to DB
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

// forwardMessageToTransactionServer forwards the received message to the Transaction Server
func forwardMessageToTransactionServer(commandID uint8, parameters commonlib.CommandParameter) (bool, error) {
	sendableCommand := commonlib.GetSendableCommand(commandID, parameters)
	response, err := commonlib.SendCommand("GET", "application/json", serverState.txPort, sendableCommand)

	log.Println("LOG: Sent command: ", sendableCommand)
	log.Println("LOG: Received response: ", response)

	if err != nil {
		log.Println("LOG: Forwarding message to Transaction Server failed with error: ", err)
		return false, err
	}
	log.Println("LOG: Forwarding message to Transaction Server succeeded.")
	return true, err
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
