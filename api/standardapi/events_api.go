// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/SummerCash/summercash-wallet-server/accounts"
)

// SetupStreams sets up all necessary event streams.
func (api *JSONHTTPAPI) SetupStreams() error {
	var accountUsernames []string // Initialize account usernames buffer

	var err error // Initialize error buffer

	accountUsernames, err = api.getAccountUsernames() // Get account usernames

	if err != nil { // Check for errors
		return err // Return found error
	}

	for _, username := range accountUsernames { // Iterate through account usernames
		api.EventServer.CreateStream(fmt.Sprintf("%s_transactions", username)) // Initialize account stream
	}

	return nil // No error occurred, return nil
}

// getAccountUsernames gets a list of account usernames.
func (api *JSONHTTPAPI) getAccountUsernames() ([]string, error) {
	var accountUsernames []string // Initialize account usernames buffer

	err := api.AccountsDatabase.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("accounts")) // Get accounts bucket

		c := bucket.Cursor() // Get cursor

		for _, accountBytes := c.First(); accountBytes != nil; _, accountBytes = c.Next() { // Iterate
			account, err := accounts.AccountFromBytes(accountBytes) // Deserialize account bytes

			if err != nil { // Check for errors
				return err // Return found error
			}

			accountUsernames = append(accountUsernames, account.Name) // Append account username to usernames list

			return nil // No error occurred, return nil
		}

		return accounts.ErrAccountDoesNotExist // Account does not exist
	}) // Read accounts

	if err != nil { // Check for errors
		return []string{}, err // Return found error
	}

	return accountUsernames, nil
}
