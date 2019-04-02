// Package accounts defines account-related helper methods and types.
// The accounts database, for example, is defined in this package.
package accounts

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/SummerCash/summercash-wallet-server/common"

	"github.com/boltdb/bolt"
	"github.com/juju/loggo"
)

var (
	// logger is the db package logger.
	logger = getDBLogger()
)

// DB is a data type representing a link to a working accounts boltdb instance.
type DB struct {
	DB *bolt.DB // DB represents the currently opened db.
}

/* BEGIN EXPORTED METHODS */

// OpenDB opens the local DB, and creates one if it doesn't already exist.
func OpenDB() (*DB, error) {
	logger.Infof("opening db instance") // Log open db

	err := common.CreateDirIfDoesNotExit(common.DBDir) // Make database directory

	if err != nil { // Check for errors
		return &DB{}, err // Return found error
	}

	db, err := bolt.Open(filepath.FromSlash(fmt.Sprintf("%s/smc_db.db", common.DBDir)), 0644, &bolt.Options{Timeout: 5 * time.Second}) // Open DB with timeout

	return &DB{
		DB: db, // Set DB
	}, nil // Return initialized db
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// getDBLogger gets the db package logger, and sets the levels of said logger.
func getDBLogger() loggo.Logger {
	logger := loggo.GetLogger("DB") // Get logger

	loggo.ConfigureLoggers("DB=INFO") // Configure loggers

	return logger // Return logger
}

/* END INTERNAL METHODS */
