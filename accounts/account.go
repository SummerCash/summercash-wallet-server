// Package accounts defines account-related helper methods and types.
// The accounts database, for example, is defined in this package.
package accounts

import "github.com/polaris-project/go-polaris/common"

// Account represents a username-password keypair linking to a private key in the accounts database.
type Account struct {
	Name string `json:"name"` // Name

	Password string `json:"password"` // Password

	Address *common.Address `json:"address"` // Address
}
