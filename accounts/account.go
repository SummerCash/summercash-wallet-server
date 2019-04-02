// Package accounts defines account-related helper methods and types.
// The accounts database, for example, is defined in this package.
package accounts

import (
	"bytes"
	"encoding/json"

	"github.com/SummerCash/go-summercash/common"
)

// Account represents a username-password keypair linking to a private key in the accounts database.
type Account struct {
	Name string `json:"name"` // Name

	PasswordHash []byte `json:"password_hash"` // Password hash

	Address common.Address `json:"address"` // Address
}

/* BEGIN EXPORTED METHODS */

// AccountFromBytes deserializes an account from a given byte array.
func AccountFromBytes(b []byte) (*Account, error) {
	account := Account{} // Init buffer

	err := json.NewDecoder(bytes.NewReader(b)).Decode(&account) // Decode into buffer

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &account, nil // No error occurred, return read value
}

// Bytes serializes a given account to a byte array.
func (account *Account) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*account, "", "  ") // Marshal

	return marshaledVal // Return marshaled val
}

/* END EXPORTED METHODS */
