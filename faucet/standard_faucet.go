// Package faucet outlines the faucet interface and its associated helper methods.
package faucet

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SummerCash/summercash-wallet-server/accounts"
	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/SummerCash/summercash-wallet-server/transactions"
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
	updatedAccount, err := faucet.AccountsDatabase.QueryAccountByUsername(account.Name) // Set to updated account

	if err != nil { // Check for errors
		return false // Cannot claim; does not exist
	}

	if updatedAccount.LastFaucetClaimTime.IsZero() { // Check not set
		return true // Has not claimed yet
	}

	if time.Now().Sub(updatedAccount.LastFaucetClaimTime).Hours() < (*faucet.GetRuleset()).GetClaimPeriod().Hours() { // Check less than 24 hours since last claim
		return false // Cannot claim
	}

	return true // Can claim
}

// AccountLastClaim gets the time at which an account last claimed from the faucet.
func (faucet *StandardFaucet) AccountLastClaim(account *accounts.Account) time.Time {
	updatedAccount, err := faucet.AccountsDatabase.QueryAccountByUsername(account.Name) // Set to updated account

	if err != nil { // Check for errors
		return time.Date(2017, time.January, 12, 0, 0, 0, 0, time.UTC) // Return despacito music video publish date
	}

	return updatedAccount.LastFaucetClaimTime // Return last claim time
}

// AmountCanClaim gets the max amount an account can claim.
// If the account has already claimed in the last 24 hours, zero is returned.
func (faucet *StandardFaucet) AmountCanClaim(account *accounts.Account) *big.Float {
	updatedAccount, err := faucet.AccountsDatabase.QueryAccountByUsername(account.Name) // Set to updated account

	if err != nil { // Check for errors
		return big.NewFloat(0)
	}

	for _, user := range faucet.Ruleset.BannedUsers() { // Iterate through banned users
		if bytes.Equal(user.Address.Bytes(), updatedAccount.Address.Bytes()) { // Check is banned user
			return big.NewFloat(0) // Return zero
		}
	}

	if time.Now().Sub(updatedAccount.LastFaucetClaimTime).Hours() < (*faucet.GetRuleset()).GetClaimPeriod().Hours() && !faucet.AccountLastClaim(updatedAccount).IsZero() { // Check less than 24 hours since last claim
		return big.NewFloat(0) // Cannot claim
	}

	return faucet.Ruleset.MaximumClaimInPeriod() // Return max claim 24 hours
}

// Claim claims a given amount from the faucet.
func (faucet *StandardFaucet) Claim(account *accounts.Account, amount *big.Float) error {
	updatedAccount, err := faucet.WorkingDB().QueryAccountByUsername(account.Name)

	if err != nil { // Check for errors
		return err // Return found error
	}

	if amountCanClaim := faucet.AmountCanClaim(updatedAccount); amountCanClaim.Cmp(amount) == -1 { // Check amount to claim less than can claim
		amountFloatVal, _ := amount.Float64() // Get float value

		canClaimFloatVal, _ := amount.Float64() // Get float value

		return fmt.Errorf("invalid claim size: %f, can claim %f", amountFloatVal, canClaimFloatVal) // Return error
	}

	err = faucet.WorkingDB().MakeFaucetClaim(updatedAccount, amount) // Make faucet claim

	if err != nil { // Check for errors
		return err // Return found error
	}

	keystoreFile, err := os.OpenFile(filepath.FromSlash(fmt.Sprintf("%s/faucet/keystore/privateKey.key", common.DataDir)), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // Open keystore dir

	if err != nil { // Check for errors
		return err // Return found error
	}

	buffer := make([]byte, 512) // Initialize pwd buffer

	_, err = keystoreFile.Read(buffer) // Read into buffer

	if err != nil { // Check for errors
		return err // Return found error
	}

	splitPassword := strings.Split(string(buffer), ":") // Split

	floatVal, _ := amount.Float64() // Get float value

	_, err = transactions.NewTransaction((*faucet).WorkingDB(), "faucet", string(splitPassword[0]+splitPassword[1]), &account.Address, floatVal, []byte("Faucet claim.")) // Initialize transaction

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// GetRuleset gets the working ruleset.
func (faucet *StandardFaucet) GetRuleset() *Ruleset {
	return &faucet.Ruleset // Return ruleset
}

/* END EXPORTED METHODS */
