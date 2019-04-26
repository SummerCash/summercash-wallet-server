// Package faucet outlines the faucet interface and its associated helper methods.
package faucet

import (
	"github.com/SummerCash/summercash-wallet-server/accounts"
	"math/big"
	"time"
)

// StandardRuleset outlines the standard ruleset implementation.
type StandardRuleset struct {
	ClaimInPeriod *big.Float `json:"claim_24hr"` // Amount can claim in 24 hours.

	ClaimPeriod time.Duration `json:"claim_period"` // Claim period

	BannedUsersList []*accounts.Account `json:"banned_users"` // Banned users.
}

/* BEGIN EXPORTED METHODS */

// NewStandardRuleset creates a new StandardRuleset instance with a given 24 hour claim amount and banned users list.
func NewStandardRuleset(claimInPeriod *big.Float, claimPeriod time.Duration, bannedUsers []*accounts.Account) *StandardRuleset {
	return &StandardRuleset{
		ClaimInPeriod:   claimInPeriod, // Set claim
		ClaimPeriod:     claimPeriod,   // Set period
		BannedUsersList: bannedUsers,   // Set banned users
	}
}

// MaximumClaimInPeriod gets the max amount one user can claim in 24 hours.
func (ruleset *StandardRuleset) MaximumClaimInPeriod() *big.Float {
	return ruleset.ClaimInPeriod // Return claim amount in 24 hours
}

// MinimumClaimInPeriod gets the min amount on user can claim in 24 hours.
func (ruleset *StandardRuleset) MinimumClaimInPeriod() *big.Float {
	return big.NewFloat(0) // No minimum for standard ruleset
}

// GetClaimPeriod gets the claim duration.
func (ruleset *StandardRuleset) GetClaimPeriod() time.Duration {
	return ruleset.ClaimPeriod // Return claim period
}

// DepositClaimCurve gets the claim curve.
func (ruleset *StandardRuleset) DepositClaimCurve() float64 {
	return 0 // No curve
}

// BannedUsers gets the list of banned users.
func (ruleset *StandardRuleset) BannedUsers() []*accounts.Account {
	return ruleset.BannedUsersList // Return banned users
}

// BanUser bans a given user.
func (ruleset *StandardRuleset) BanUser(user *accounts.Account) {
	(*ruleset).BannedUsersList = append((*ruleset).BannedUsersList, user) // Ban user
}

/* END EXPORTED METHODS */
