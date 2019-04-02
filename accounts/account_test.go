// Package accounts defines account-related helper methods and types.
// The accounts database, for example, is defined in this package.
package accounts

import (
	"bytes"
	"testing"

	"github.com/SummerCash/go-summercash/common"

	"github.com/SummerCash/summercash-wallet-server/crypto"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestBytesAccount tests the functionality of the Bytes() helper method.
func TestBytesAccount(t *testing.T) {
	address, err := common.StringToAddress("0x040028d536d5351e83fbbec320c194629ace") // Get addr value

	if err != nil { // Check for errors
		t.Error(err) // Log found error
		t.FailNow()  // Panic
	}

	account := &Account{
		Name:         "Dowland Aiello",                               // Set name
		PasswordHash: crypto.Salt([]byte("despacito_despacito_lol")), // Set password
		Address:      address,                                        // Set address
	}

	deserializedAccount, err := AccountFromBytes(account.Bytes()) // Deserialize account bytes

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if !bytes.Equal(deserializedAccount.Address.Bytes(), address.Bytes()) { // Check address not matching
		t.Fatal("address deserialized incorrectly") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
