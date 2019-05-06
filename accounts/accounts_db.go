// Package accounts defines account-related helper methods and types.
// The accounts database, for example, is defined in this package.
package accounts

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	rand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/juju/loggo"

	"github.com/SummerCash/go-summercash/accounts"
	summercashCommon "github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/types"
	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/SummerCash/summercash-wallet-server/crypto"
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

	database, err := bolt.Open(filepath.FromSlash(fmt.Sprintf("%s/smc_db.db", common.DBDir)), 0644, &bolt.Options{Timeout: 5 * time.Second}) // Open DB with timeout

	if err != nil { // Check for errors
		return &DB{}, err // Return found error
	}

	db := &DB{
		DB: database, // Set DB
	} // Initialize DB

	err = db.DB.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(accountsBucket) == nil { // Check first account
			return errors.New("must make faucet account") // Return must make error
		}

		return nil // No error occurred, return nil
	}) // Create faucet account

	if err != nil && err.Error() == "must make faucet account" { // Check must make faucet account
		var privateKey *ecdsa.PrivateKey // Initialize private key buffer

		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

		if err != nil { // Check for errors
			return &DB{}, err // Return found error
		}

		err = common.CreateDirIfDoesNotExit(fmt.Sprintf("%s/faucet/keystore", common.DataDir)) // Create faucet keystore dir

		if err != nil { // Check for errors
			return &DB{}, err // Return found error
		}

		keystoreFile, err := os.OpenFile(filepath.FromSlash(fmt.Sprintf("%s/faucet/keystore/privateKey.key", common.DataDir)), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // Open keystore dir

		if err != nil { // Check for errors
			return &DB{}, err // Return found error
		}

		defer keystoreFile.Close() // Close keystore file

		_, err = keystoreFile.WriteString(privateKey.X.String() + ":" + privateKey.Y.String()) // Write pwd

		if err != nil { // Check for errors
			return &DB{}, err // Return found error
		}

		_, err = db.CreateNewAccount("faucet", privateKey.X.String()+privateKey.Y.String()) // Create faucet account

		if err != nil { // Check for errors
			return &DB{}, err // Return found error
		}
	}

	if err != nil { // Check for errors
		return &DB{}, err // Return found error
	}

	return db, nil // Return initialized db
}

// CloseDB closes the db.
func (db *DB) CloseDB() error {
	return db.DB.Close() // Close db
}

// AddNewAccount adds a new account to the list of accounts in the working database.
func (db *DB) AddNewAccount(name string, password string, address string) (*Account, error) {
	parsedAddress, err := summercashCommon.StringToAddress(address) // Parse hex address

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	account := &Account{
		Name:         name,                          // Set name
		PasswordHash: crypto.Salt([]byte(password)), // Set password hash
		Address:      parsedAddress,                 // Set address
	}

	err = db.CreateAccountsBucketIfNotExist() // Create accounts bucket

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	err = db.DB.Update(func(tx *bolt.Tx) error {
		accountsBucket := tx.Bucket(accountsBucket) // Get accounts bucket

		if alreadyExists := accountsBucket.Get(crypto.Sha3([]byte(name))); alreadyExists != nil { // Check already exists
			return ErrAccountAlreadyExists // Return error
		}

		return accountsBucket.Put(crypto.Sha3([]byte(name)), account.Bytes()) // Put account
	}) // Add new account to DB

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	return account, nil // Return address
}

// IssueAccountToken issues a new account token.
func (db *DB) IssueAccountToken(username, password string) (string, error) {
	account, err := db.QueryAccountByUsername(username) // Query account

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	if !db.Auth(username, password) { // Check should not issue token
		return "", ErrPasswordInvalid // Invalid
	}

	token, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate key

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	(*account).Tokens = append(account.Tokens, string(crypto.Sha3(append(token.X.Bytes(), token.Y.Bytes()...)))) // Add token to account

	err = db.DB.Update(func(tx *bolt.Tx) error {
		accountsBucket := tx.Bucket(accountsBucket) // Get accounts bucket

		return accountsBucket.Put(crypto.Sha3([]byte(account.Name)), account.Bytes()) // Update account
	})

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return string(crypto.Sha3(append(token.X.Bytes(), token.Y.Bytes()...))), nil // Return token
}

// ValidateAccountToken checks whether or not a given token is valid.
func (db *DB) ValidateAccountToken(account *Account, token string) bool {
	for _, currentToken := range account.Tokens { // Iterate through account tokens
		if currentToken == token { // Check token matches
			return true // Valid
		}
	}

	return false // Invalid token
}

// MakeFaucetClaim makes a faucet claim for a given account.
func (db *DB) MakeFaucetClaim(account *Account, amount *big.Float) error {
	err := db.CreateAccountsBucketIfNotExist() // Create accounts bucket

	if err != nil { // Check for errors
		return err // Return found error
	}

	return db.DB.Update(func(tx *bolt.Tx) error {
		accountsBucket := tx.Bucket(accountsBucket) // Get accounts bucket

		(*account).LastFaucetClaimTime = time.Now() // Set last claim time
		(*account).LastFaucetClaimAmount = amount   // Set claim amount

		return accountsBucket.Put(crypto.Sha3([]byte(account.Name)), account.Bytes()) // Put account
	}) // Add new account to DB
}

// GetUserBalance calculates the balance of a particular account.
func (db *DB) GetUserBalance(username string) (*big.Float, error) {
	account, err := db.QueryAccountByUsername(username) // Query account

	if err != nil { // Check for errors
		return big.NewFloat(0), err // Return found error
	}

	chain, err := types.ReadChainFromMemory(account.Address) // Read account

	if err != nil { // Check for errors
		return big.NewFloat(0), err // Return found error
	}

	return chain.CalculateBalance(), nil // Return calculated balance
}

// GetUserTransactions fetches the list of transactions for a particular account.
func (db *DB) GetUserTransactions(username string) ([]*types.Transaction, error) {
	account, err := db.QueryAccountByUsername(username) // Query account

	if err != nil { // Check for errors
		return []*types.Transaction{}, err // Return found error
	}

	chain, err := types.ReadChainFromMemory(account.Address) // Read account

	if err != nil { // Check for errors
		return []*types.Transaction{}, err // Return found error
	}

	return chain.Transactions, nil // Return account chain transactions
}

// CreateNewAccount creates a new account with a given name and password.
// Returns the new account's address and an error (if applicable).
func (db *DB) CreateNewAccount(name string, password string) (*Account, error) {
	account, err := accounts.NewAccount() // Create new account

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	err = account.WriteToMemory() // Write account to persistent memory

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	accountInstance := &Account{
		Name:         name,                          // Set name
		PasswordHash: crypto.Salt([]byte(password)), // Set password hash
		Address:      account.Address,               // Set address
	}

	err = db.CreateAccountsBucketIfNotExist() // Create accounts bucket

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	err = db.DB.Update(func(tx *bolt.Tx) error {
		accountsBucket := tx.Bucket(accountsBucket) // Get accounts bucket

		if alreadyExists := accountsBucket.Get(crypto.Sha3([]byte(name))); alreadyExists != nil { // Check already exists
			return ErrAccountAlreadyExists // Return error
		}

		return accountsBucket.Put(crypto.Sha3([]byte(name)), accountInstance.Bytes()) // Put account
	}) // Add new account to DB

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	return accountInstance, nil // Return address
}

// Auth checks that a given user can be authenticated.
func (db *DB) Auth(name string, password string) bool {
	account, err := db.QueryAccountByUsername(name) // Query by username

	if err != nil { // Check for errors
		return false // Return not valid
	}

	return crypto.VerifySalted(account.PasswordHash, password) || db.ValidateAccountToken(account, password) // Verify salt / token
}

// DeleteAccount deletes an account from the working DB.
func (db *DB) DeleteAccount(name string, password string) error {
	err := db.CreateAccountsBucketIfNotExist() // Create accounts bucket

	if err != nil { // Check for errors
		return err // Return found error
	}

	if !db.Auth(name, password) { // Auth
		return ErrPasswordInvalid // Return error
	}

	return db.DB.Update(func(tx *bolt.Tx) error {
		accountsBucket := tx.Bucket(accountsBucket) // Get accounts bucket

		return accountsBucket.Delete(crypto.Sha3([]byte(name))) // Delete account
	}) // Update account info
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

	err = db.CreateAccountsBucketIfNotExist() // Create accounts bucket

	if err != nil { // Check for errors
		return err // Return found error
	}

	return db.DB.Update(func(tx *bolt.Tx) error {
		accountsBucket := tx.Bucket(accountsBucket) // Get accounts bucket

		return accountsBucket.Put(crypto.Sha3([]byte(name)), account.Bytes()) // Put account
	}) // Update account info
}

// QueryAccountByUsername queries the database for an account with a given username.
func (db *DB) QueryAccountByUsername(name string) (*Account, error) {
	var accountBuffer *Account // Initialize account buffer

	err := db.CreateAccountsBucketIfNotExist() // Create accounts bucket

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	err = db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(accountsBucket) // Get accounts bucket

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

// QueryAccountByAddress queries the database for an account with a given address.
func (db *DB) QueryAccountByAddress(address summercashCommon.Address) (*Account, error) {
	var accountBuffer *Account // Initialize account buffer

	err := db.CreateAccountsBucketIfNotExist() // Create accounts bucket

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	err = db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(accountsBucket) // Get accounts bucket

		c := bucket.Cursor() // Get cursor

		for _, accountBytes := c.First(); accountBytes != nil; _, accountBytes = c.Next() { // Iterate
			account, err := AccountFromBytes(accountBytes) // Deserialize account bytes

			if bytes.Equal(account.Address.Bytes(), address.Bytes()) && err == nil { // Check addresses equivalent
				accountBuffer = account // Set account

				return nil // No error occurred, return nil
			}
		}

		return ErrAccountDoesNotExist // Account does not exist
	}) // Read account

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	return accountBuffer, nil // Return read account
}

// CreateAccountsBucketIfNotExist creates the accounts bucket if it doesn't already exist.
func (db *DB) CreateAccountsBucketIfNotExist() error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(accountsBucket) // Create bucket

		return err // Return error
	}) // Create bucket
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
