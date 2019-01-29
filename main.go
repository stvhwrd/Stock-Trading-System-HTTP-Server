package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	commonlib "../SENG-468-Common-Lib"
)

// ServerState holds information about the server's state and configuration
type ServerState struct {
	dbPort   int
	tsPort   int
	httpPort int
}

var serverState = ServerState{}

func init() {
	// Parse and process CLI flags
	flag.IntVar(&serverState.dbPort, "dbport", -1, "[REQUIRED] the port on which the DATABASE server is running, eg. --dbport=8080")
	flag.IntVar(&serverState.httpPort, "httpport", 8084, "[optional -- default is port 8084] the port on which *this* HTTP server is running, eg. --httpport=8084")
	flag.IntVar(&serverState.tsPort, "tsport", -1, "[REQUIRED] the port on which the TRANSACTION server is running, eg. --tsport=8082")
	flag.Parse()

	// Force flags as required
	if serverState.dbPort == -1 || serverState.tsPort == -1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	// Display server state values
	fmt.Printf("%+v\n", serverState)

	// Fire up the server
	// portNumString := ":" + strconv.Itoa(serverState.httpPort)
	commonlib.StartServer(serverState.httpPort, requestHandler)

	for {
		// Keep the fun running
	}
}

func requestHandler(w http.ResponseWriter, response *http.Request) {
	fmt.Println("Message received:", response.Method)

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("An error occurred with this message:", err)
		fmt.Println("---------- End Message ----------")
		log.Fatal(err)
		return
	}

	defer response.Body.Close()
	fmt.Println("Message body:", string(responseBody))
	// fmt.Println("Message command:", commonlib.GetMessageType(responseBody))
	defer fmt.Println("---------- End Message ----------")

	fmt.Println("Parsing message...")

	// Get the command and its parameters back from the response.
	command, parameters := commonlib.GetCommandFromMessage(responseBody)
	// fmt.Println("Command:", command)
	// fmt.Println("CommandParameters:", parameters)
	processCommands(w, responseBody, command, parameters)
}

func processCommands(responseWriter http.ResponseWriter, originalResponse []byte, command uint8, parameters commonlib.CommandParameter) {
	shouldReply := true
	var reply string
	var err error

	switch command {
	case commonlib.AddCommand:
		if amount, _ := strconv.ParseFloat(parameters.Amount, 64); amount < 0 {
			fmt.Println("Cannot add negative Amount")
		} else {
			fmt.Println("Routed AddCommand to Transaction server.")
			// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
		}
	case commonlib.QuoteCommand:
		fmt.Println("Routed QuoteCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.BuyCommand:
		fmt.Println("Routed BuyCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.CommitBuyCommand:
		fmt.Println("Routed CommitBuyCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.CancelBuyCommand:
		fmt.Println("Routed CancelBuyCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.SellCommand:
		fmt.Println("Routed SellCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.CommitSellCommand:
		fmt.Println("Routed CommitSellCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.CancelSellCommand:
		fmt.Println("Routed CancelSellCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.SetBuyAmountCommand:
		fmt.Println("Routed SetBuyAmountCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.CancelSetBuyCommand:
		fmt.Println("Routed CancelSetBuyCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.SetBuyTriggerCommand:
		fmt.Println("Routed SetBuyTriggerCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.SetSellAmountCommand:
		fmt.Println("Routed SetSellAmountCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.SetSellTriggerCommand:
		fmt.Println("Routed SetSellTriggerCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.CancelSetSellCommand:
		fmt.Println("Routed CancelSetSellCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.DumplogCommand:
		fmt.Println("Routed DumplogCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.DumplogAllCommand:
		fmt.Println("Routed DumplogAllCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	case commonlib.DisplaySummaryCommand:
		fmt.Println("Routed DisplaySummaryCommand to Transaction server.")
		// reply, err = commonlib.SendCommand("POST", "application/json", serverState.tsPort, originalResponse)
	default:
		fmt.Println("You sent the wrong command id who are you trying to fool MTFK")
	}

	// If an error occurred sending a command and getting a reply, log details.
	if err != nil {
		log.Fatal("Failed to send command from transaction server:", command, parameters, ". Error:", err)
		reply = err.Error()
	}

	if shouldReply {
		fmt.Println("Replying back.")
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.Write([]byte(reply))
	}

}
