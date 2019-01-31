package main

import commonlib "github.com/kurtd5105/SENG-468-Common-Lib"

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
