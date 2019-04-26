// Package faucet outlines the faucet interface and its associated helper methods.
package faucet

import (
	"github.com/SummerCash/summercash-wallet-server/accounts"

	"math/big"
	"time"
)

// Ruleset defines all methods necessary to implement a ruleset.
type Ruleset interface {
	MaximumClaimInPeriod() *big.Float // Get max amount can claim in 24 hours.
	MinimumClaimInPeriod() *big.Float // Get min amount can claim in 24 hours.

	GetClaimPeriod() time.Duration // Get duration between possible claims

	DepositClaimCurve() float64 // Get amount to multiply possible claim by deposit by.

	BannedUsers() []*accounts.Account // Get banned users.
	BanUser(*accounts.Account)        // Ban user.
}
