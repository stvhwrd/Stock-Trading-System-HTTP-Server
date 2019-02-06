package main

import (
	"log"
	"strconv"

	"github.com/kurtd5105/SENG-468-Common-Lib"
)

var transactionCommandStringMap = map[string]int{
	"ADD":              commonlib.AddCommand,
	"BUY":              commonlib.BuyCommand,
	"CANCEL_BUY":       commonlib.CancelBuyCommand,
	"CANCEL_SELL":      commonlib.CancelSellCommand,
	"CANCEL_SET_BUY":   commonlib.CancelSetBuyCommand,
	"CANCEL_SET_SELL":  commonlib.CancelSetSellCommand,
	"COMMIT_BUY":       commonlib.CommitBuyCommand,
	"COMMIT_SELL":      commonlib.CommitSellCommand,
	"DISPLAY_SUMMARY":  commonlib.DisplaySummaryCommand,
	"DUMPLOG":          commonlib.DumplogCommand,
	"QUOTE":            commonlib.QuoteCommand,
	"SELL":             commonlib.SellCommand,
	"SET_BUY_AMOUNT":   commonlib.SetBuyAmountCommand,
	"SET_BUY_TRIGGER":  commonlib.SetBuyTriggerCommand,
	"SET_SELL_AMOUNT":  commonlib.SetSellAmountCommand,
	"SET_SELL_TRIGGER": commonlib.SetSellTriggerCommand}

func buildAndSendMessage(payload JSONPayload) {
	var commandID int
	var parameters = commonlib.CommandParameter{}
	var destinationIP string
	var destinationPort int

	// HACK: this is messy but there are only three commands that need to be handled individually.
	// DUMPLOGs go to the logging server,
	// DISPLAY_SUMMARY goes to the database,
	// the rest go to the transaction server.
	switch payload.Command {
	case "DUMPLOG":
		if len(payload.UserID) > 1 {
			commandID = commonlib.DumplogCommand
			parameters.UserID = payload.UserID
			parameters.Filename = payload.Filename
		} else {
			commandID = commonlib.DumplogAllCommand
			parameters.Filename = payload.Filename
		}

	case "DISPLAY_SUMMARY":
		commandID = commonlib.DisplaySummaryCommand
		parameters.UserID = payload.UserID

	default:
		commandID = transactionCommandStringMap[payload.Command]
		destinationIP = serverConfig.transaction.ipAddress
		destinationPort = serverConfig.transaction.port
		parameters.UserID = payload.UserID
		parameters.Amount = payload.Amount
		parameters.Filename = payload.Filename
		parameters.StockSymbol = payload.StockSymbol
	}

	sendMessage(destinationIP, destinationPort, payload.Command, commandID, parameters)
}

func sendMessage(destinationIP string, destinationPort int, originalCommand string, commandID int, parameters commonlib.CommandParameter) {
	log.Printf("Sending POST request with command to %s:%d/ \n\n", destinationIP, destinationPort)
	log.Printf("\tCOMMAND: %s (#%d)\n", originalCommand, commandID)
	log.Printf("\tPARAMETERS: %+v\n\n", parameters)

	response, err := commonlib.SendCommand("POST", "application/json", destinationIP+":"+strconv.Itoa(destinationPort), commonlib.GetSendableCommand(uint8(commandID), parameters))
	if err != nil {
		log.Printf("-- Error sending command --")
		panic(err)
	}
	log.Printf("Received response:\n%s\n", response)
}
