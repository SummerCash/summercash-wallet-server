// Package faucet outlines the faucet interface and its associated helper methods.
package faucet

import (
	"math/big"
)

// Ruleset defines all methods necessary to implement a ruleset.
type Ruleset interface {
	MaximumClaim24hr() *big.Float  // Get max amount can claim in 24 hours.
	MininmumClaim24hr() *big.Float // Get min amount can claim in 24 hours.

	DepositPerClaim(claim *big.Float) // Get deposit required at claim amount.
}
