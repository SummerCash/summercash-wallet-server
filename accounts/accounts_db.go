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

	// ErrAccountDoesNotExist is an error definition describing an account value of nil.
	ErrAccountDoesNotExist = errors.New("no account exists with the given username")

	// ErrPasswordInvalid is an error definition describing an invalid password value.
	ErrPasswordInvalid = errors.New("invalid password")
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
func (db *DB) CreateNewAccount(name string, password string) (string, error) {
	account, err := accounts.NewAccount() // Create new account

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	accountInstance := &Account{
		Name:         name,                          // Set name
		PasswordHash: crypto.Salt([]byte(password)), // Set password hash
		Address:      account.Address,               // Set address
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

// ResetAccountPassword resets an accounts password.
func (db *DB) ResetAccountPassword(name string, oldPassword string, newPassword string) error {
	account, err := db.QueryAccountByUsername(name) // Query by username

	if err != nil { // Check for errors
		return err // Return found error
	}

	valid := crypto.VerifySalted(account.PasswordHash, oldPassword) // Verify salt

	if valid != true { // Check for errors
		return ErrPasswordInvalid // Return found error
	}

	(*account).PasswordHash = crypto.Salt([]byte(newPassword)) // Set salt

	return db.DB.Update(func(tx *bolt.Tx) error {
		accountsBucket, err := tx.CreateBucketIfNotExists(accountsBucket) // Create accounts bucket

		if err != nil { // Check for errors
			return err // Return found error
		}

		return accountsBucket.Put(crypto.Sha3([]byte(name)), account.Bytes()) // Put account
	}) // Update account info
}

// QueryAccountByUsername queries the database for an account with a given username.
func (db *DB) QueryAccountByUsername(name string) (*Account, error) {
	var accountBuffer *Account // Initialize account buffer

	err := db.DB.View(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(accountsBucket) // Get account bucket

		if err != nil { // Check for errors
			return err // Return found error
		}

		accountBytes := bucket.Get(crypto.Sha3([]byte(name))) // Get account at hash

		if accountBytes == nil { // Check no account at hash
			return ErrAccountDoesNotExist // Return error
		}

		accountBuffer, err = AccountFromBytes(accountBytes) // Deserialize account bytes

		return err // Return error
	}) // Read account

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	return accountBuffer, nil // Return read account
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