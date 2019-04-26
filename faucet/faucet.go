// Package faucet outlines the faucet interface and its associated helper methods.
package faucet

import (
	"github.com/SummerCash/summercash-wallet-server/accounts"
	"math/big"
	"time"
)

// Faucet defines all methods that should be associated with a faucet.
type Faucet interface {
	WorkingDB() *accounts.DB // Get Accounts DB ref.

	AccountCanClaim(account *accounts.Account) bool       // Check account can claim summercash.
	AccountLastClaim(account *accounts.Account) time.Time // Get last time account claimed.

	AmountCanClaim(account *accounts.Account) *big.Float // Get the max amount an account can claim.

	Claim(account *accounts.Account, amount *big.Float) error // Claim an amount of summercash from the faucet.

	GetRuleset() *Ruleset // Get ruleset
}
