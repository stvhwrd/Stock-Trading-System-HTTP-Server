package main

import (
	"log"
	"os"
)

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
