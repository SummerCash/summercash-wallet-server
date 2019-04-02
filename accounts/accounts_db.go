// Package accounts defines account-related helper methods and types.
// The accounts database, for example, is defined in this package.
package accounts

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/SummerCash/go-summercash/accounts"
	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/SummerCash/summercash-wallet-server/crypto"

	"github.com/boltdb/bolt"
	"github.com/juju/loggo"
)

var (
	// ErrAccountAlreadyExists is an error definition describing an attempted duplicate put.
	ErrAccountAlreadyExists = errors.New("account already exists")
)

var (
	// accountsBucket is the accounts bucket key definition.
	accountsBucket = []byte("accounts")

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

// CreateNewAccount creates a new account with a given name and password.
// Returns the new account's address and an error (if applicable).
func (db *DB) CreateNewAccount(name string, passwordHash []byte) (string, error) {
	account, err := accounts.NewAccount() // Create new account

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	accountInstance := &Account{
		Name:         name,            // Set name
		PasswordHash: passwordHash,    // Set password hash
		Address:      account.Address, // Set address
	}

	err = db.DB.Update(func(tx *bolt.Tx) error {
		accountsBucket, err := tx.CreateBucketIfNotExists(accountsBucket) // Create accounts bucket

		if err != nil { // Check for errors
			return err // Return found error
		}

		if alreadyExists := accountsBucket.Get(crypto.Sha3([]byte(name))); alreadyExists != nil { // Check already exists
			return ErrAccountAlreadyExists // Return error
		}

		return accountsBucket.Put(crypto.Sha3([]byte(name)), accountInstance.Bytes()) // Put account
	}) // Add new account to DB

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return account.Address.String(), nil // Return address
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
