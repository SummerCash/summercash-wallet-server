// Package faucet outlines the faucet interface and its associated helper methods.
package faucet

import (
	"github.com/SummerCash/summercash-wallet-server/accounts"
	"time"
)

// StandardFaucet outlines a standard faucet conforming to the standard ruleset.
type StandardFaucet struct {
	Ruleset Ruleset // Faucet ruleset

	AccountsDatabase *accounts.DB // Accounts database
}

/* BEGIN EXPORTED METHODS */

// NewStandardFaucet initializes a new standard faucet.
func NewStandardFaucet(ruleset Ruleset, accountsDB *accounts.DB) *StandardFaucet {
	return &StandardFaucet{
		Ruleset:          ruleset,    // Set ruleset
		AccountsDatabase: accountsDB, // Set DB
	} // Return new faucet
}

// WorkingDB gets a reference to the faucet working database.
func (faucet *StandardFaucet) WorkingDB() *accounts.DB {
	return faucet.AccountsDatabase // Return DB
}

// AccountCanClaim checks if a given account can claim summercash.
func (faucet *StandardFaucet) AccountCanClaim(account *accounts.Account) bool {
	if account.LastFaucetClaimTime.Sub(time.Now()).Hours() < 24 { // Check less than 24 hours
		return false // Cannot claim
	}

	return true // Can claim
}

/* END EXPORTED METHODS */
