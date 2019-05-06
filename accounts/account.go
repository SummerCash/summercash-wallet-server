// Package accounts defines account-related helper methods and types.
// The accounts database, for example, is defined in this package.
package accounts

import (
	"bytes"
	"encoding/json"
	"math/big"
	"time"

	"github.com/SummerCash/go-summercash/common"
)

// Account represents a username-password keypair linking to a private key in the accounts database.
type Account struct {
	Name string `json:"name"` // Name

	PasswordHash []byte `json:"password_hash"` // Password hash

	Address common.Address `json:"address"` // Address

	LastFaucetClaimTime   time.Time  `json:"last_claim_time"`   // Last claim time
	LastFaucetClaimAmount *big.Float `json:"last_claim_amount"` // Last claim amount

	Tokens []string `json:"tokens"` // Account tokens
}

// jsonAccount represents a JSON-friendly account.
type jsonAccount struct {
	Name string `json:"name"` // Name

	PasswordHash []byte `json:"password_hash"` // Password hash

	HexAddress string `json:"address"` // Address
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

// String serializes a given account to a JSON string.
func (account *Account) String() string {
	jsonAccount := jsonAccount{
		Name:         account.Name,             // Set name
		PasswordHash: account.PasswordHash,     // Set password hash
		HexAddress:   account.Address.String(), // Set hex address
	} // Initialize JSON account instance

	marshaledVal, _ := json.MarshalIndent(jsonAccount, "", "  ") // Marshal

	return string(marshaledVal) // Return marshaled val
}

// Bytes serializes a given account to a byte array.
func (account *Account) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*account, "", "  ") // Marshal

	return marshaledVal // Return marshaled val
}

/* END EXPORTED METHODS */
