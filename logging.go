package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	commonlib "github.com/kurtd5105/SENG-468-Common-Lib"
)

// buildLog builds a valid LogCommand
func buildLog(logMessage string, logType int, logParameters commonlib.LogCommandParameter) []byte {

	logParameters.Timestamp = getUnixTimestampAsString()
	logParameters.Server = "Web"
	var logCommandID uint8

	log.Println(logMessage)

	switch logType {

	// Required fields: Timestamp, Server, TransactionNum, Command
	// Optional fields: Username, StockSymbol, Filename, Funds, DebugMessage
	case commonlib.DebugType:
		logCommandID = uint8(commonlib.DebugType)
		logParameters.DebugMessage = logMessage

	// Required fields: Timestamp, Server, TransactionNum, Command
	// Optional fields: Username, StockSymbol, Filename, Funds, ErrorMessage
	case commonlib.ErrorEventType:
		logCommandID = uint8(commonlib.ErrorEventType)
		logParameters.ErrorMessage = logMessage

	// Required fields: Timestamp, Server, TransactionNum, Command
	// Optional fields: Username, StockSymbol, Filename, Funds
	case commonlib.SystemEventType:
		logCommandID = uint8(commonlib.SystemEventType)

	// Required fields: Timestamp, Server, TransactionNum, Command
	// Optional fields: Username, StockSymbol, Filename, Funds
	case commonlib.UserCommandType:
		logCommandID = uint8(commonlib.UserCommandType)

	default:
		log.Println("Web server not set up to handle this log type: " + string(logType))
	}

	sendableLogCommand := commonlib.GetSendableLogCommand(logCommandID, logParameters)
	if len(sendableLogCommand) == 0 {
		fmt.Println("GetSendableLogCommand returned a 0-length response")
	}
	return sendableLogCommand
}

// sendLog sends the given log as a message to the logging server
func sendLog(sendableLogCommand []byte) {
	replyBody, err := commonlib.SendCommand("POST", "application/json", state.loggingServerAddressAndPort, sendableLogCommand)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	log.Printf(replyBody)
}

// getUnixTimestampAsString returns a string representing the number of milliseconds (UTC) since the UNIX epoch
func getUnixTimestampAsString() string {
	// Format required is milliseconds
	milliseconds := time.Now().UTC().UnixNano() / int64(1000)
	return strconv.FormatInt(milliseconds, 10)
}
