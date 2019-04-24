// Package faucet outlines the faucet interface and its associated helper methods.
package faucet

import (
	"github.com/SummerCash/summercash-wallet-server/accounts"
	"math/big"
)

// StandardRuleset outlines the standard ruleset implementation.
type StandardRuleset struct {
	Claim24hr *big.Float `json:"claim_24hr"` // Amount can claim in 24 hours.

	BannedUsersList []*accounts.Account `json:"banned_users"` // Banned users.
}

/* BEGIN EXPORTED METHODS */

// NewStandardRuleset creates a new StandardRuleset instance with a given 24 hour claim amount and banned users list.
func NewStandardRuleset(claim24hr *big.Float, bannedUsers []*accounts.Account) *StandardRuleset {
	return &StandardRuleset{
		Claim24hr:       claim24hr,   // Set claim
		BannedUsersList: bannedUsers, // Set banned users
	}
}

// MaximumClaim24hr gets the max amount one user can claim in 24 hours.
func (ruleset *StandardRuleset) MaximumClaim24hr() *big.Float {
	return ruleset.Claim24hr // Return claim amount in 24 hours
}

// MinimumClaim24hr gets the min amount on user can claim in 24 hours.
func (ruleset *StandardRuleset) MinimumClaim24hr() *big.Float {
	return big.NewFloat(0) // No minimum for standard ruleset
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
