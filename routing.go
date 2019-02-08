package main

import (
	commonlib "github.com/kurtd5105/SENG-468-Common-Lib"
)

func getCommandID(command string, userID string) uint8 {
	// This map contains all possible commands except DUMPLOGs since those will be handled individually
	commandStringMap := map[string]int{
		"ADD":              commonlib.AddCommand,
		"BUY":              commonlib.BuyCommand,
		"CANCEL_BUY":       commonlib.CancelBuyCommand,
		"CANCEL_SELL":      commonlib.CancelSellCommand,
		"CANCEL_SET_BUY":   commonlib.CancelSetBuyCommand,
		"CANCEL_SET_SELL":  commonlib.CancelSetSellCommand,
		"COMMIT_BUY":       commonlib.CommitBuyCommand,
		"COMMIT_SELL":      commonlib.CommitSellCommand,
		"DISPLAY_SUMMARY":  commonlib.DisplaySummaryCommand,
		"QUOTE":            commonlib.QuoteCommand,
		"SELL":             commonlib.SellCommand,
		"SET_BUY_AMOUNT":   commonlib.SetBuyAmountCommand,
		"SET_BUY_TRIGGER":  commonlib.SetBuyTriggerCommand,
		"SET_SELL_AMOUNT":  commonlib.SetSellAmountCommand,
		"SET_SELL_TRIGGER": commonlib.SetSellTriggerCommand}

	if command == "DUMPLOG" {
		if userID =! "" {
			return commonlib.DumplogCommand
		}
		return commonlib.DumplogAllCommand
	}
	return uint8(commandStringMap[command])
}

func getDestinationServer(commandID uint8) string {
	switch commandID {
	// TODO: Display Summary command will need to be handled in isolation since
	// it requires information from both logging server and accounts database
	case commonlib.DisplaySummaryCommand:
		return state.databaseServerAddressAndPort
	case commonlib.DumplogAllCommand, commonlib.DumplogCommand:
		return state.loggingServerAddressAndPort
	default:
		return state.transactionServerAddressAndPort
	}
}
