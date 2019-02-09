package main

import (
	"fmt"

	commonlib "github.com/kurtd5105/SENG-468-Common-Lib"
)

// "timestamp":** `int`  // UNIX timestamp
// "server":** `string`  // name of server eg. "Transaction Server"
// "transactionNum":** `int`  // enumeration of transaction
// "command":** `string`  // one of the known commands ("BUY", "SELL" etc)
// "username":** `string`  // **optional** alphanumeric username
// "stockSymbol":** `string`  // **optional** three-char alphanumeric stock symbol ("NWC", "BTK")
// "filename":** `string`  // **optional** name of file to be written to
// "funds":** `decimal`  // **optional** dollars and cents (eg 24.99)
// "debugMessage":** `string`  // **optional** debug message or description
//
// DebugType
func buildLogDebug() (uint8, commonlib.LogCommandParameter) {
	logParameters := commonlib.LogCommandParameter{
		Server:         "Web",
		TransactionNum: "0001",
		Command:        "QUOTE",
	}

	commandID := uint8(commonlib.DebugType)

	return commandID, logParameters
}

// "timestamp":** `int`  // UNIX timestamp
// "server":** `string`  // name of server eg. "Transaction Server"
// "transactionNum":** `int`  // enumeration of transaction
// "command":** `string`  // one of the known commands ("BUY", "SELL" etc)
// "username":** `string`  // **optional** alphanumeric username
// "stockSymbol":** `string`  // **optional** three-char alphanumeric stock symbol ("NWC", "BTK")
// "filename":** `string`  // **optional** name of file to be written to
// "funds":** `decimal`  // **optional** dollars and cents (eg 24.99)
// "errorMessage":** `string`  // **optional** error message or description
//
// ErrorEventType
func buildLogErrorEvent() (uint8, commonlib.LogCommandParameter) {
	logParameters := commonlib.LogCommandParameter{
		Server:         "Web",
		TransactionNum: "0002",
		Command:        "SET_BUY_TRIGGER",
	}

	logCommandID := uint8(commonlib.ErrorEventType)

	return logCommandID, logParameters
}

// "timestamp":** `int`  // UNIX timestamp
// "server":** `string`  // name of server eg. "Transaction Server"
// "transactionNum":** `int`  // enumeration of transaction
// "command":** `string`  // one of the known commands ("BUY", "SELL" etc)
// "username":** `string`  // **optional** alphanumeric username
// "stockSymbol":** `string`  // **optional** three-char alphanumeric stock symbol ("NWC", "BTK")
// "filename":** `string`  // **optional** name of file to be written to
// "funds":** `decimal`  // **optional** dollars and cents (eg 24.99)
//
// SystemEventType
func buildLogSystemEvent() (uint8, commonlib.LogCommandParameter) {
	logParameters := commonlib.LogCommandParameter{
		Server:         "Web",
		TransactionNum: "0003",
		Command:        "DUMPLOG",
		Username:       "debugTestUser",
		LogStockSymbol: "NWC",
	}

	logCommandID := uint8(commonlib.SystemEventType)

	return logCommandID, logParameters
}

// "timestamp":** `int`  // UNIX timestamp
// "server":** `string`  // name of server eg. "Transaction Server"
// "transactionNum":** `int`  // enumeration of transaction
// "command":** `string`  // one of the known commands ("BUY", "SELL" etc)
// "username":** `string`  // **optional** alphanumeric username
// "stockSymbol":** `string`  // **optional** three-char alphanumeric stock symbol ("NWC", "BTK")
// "filename":** `string`  // **optional** name of file to be written to
// "funds":** `decimal`  // **optional** dollars and cents (eg 24.99)
//
// UserCommandType
func buildLogUserCommand() (uint8, commonlib.LogCommandParameter) {
	logParameters := commonlib.LogCommandParameter{
		Server:         "Web",
		TransactionNum: "0004",
		Command:        "QUOTE",
		Username:       "debugTestUser",
		LogStockSymbol: "NWC",
	}

	logCommandID := uint8(commonlib.UserCommandType)

	return logCommandID, logParameters
}

// sendLog sends the given log as a message to the logging server
func sendLog(logCommandID uint8, logParameters commonlib.LogCommandParameter) {
	request := commonlib.GetSendableLogCommand(logCommandID, logParameters)
	if len(request) == 0 {
		fmt.Println("GetSendableLogCommand returned a 0-length byte slice")
	}

	replyBody, err := commonlib.SendCommand("POST", "application/json", state.loggingServerAddressAndPort, request)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Printf(replyBody)
}
