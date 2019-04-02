// Package accounts defines account-related helper methods and types.
// The accounts database, for example, is defined in this package.
package accounts

import (
	"testing"

	"github.com/juju/loggo"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestOpenDB tests the functionality of the OpenDB() helper method.
func TestOpenDB(t *testing.T) {
	_, err := OpenDB() // Open db

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */

/* BEGIN INTERNAL METHODS TESTS */

// TestGetDBLogger tests the functionality of the GetDBLogger() helper method.
func TestGetDBLogger(t *testing.T) {
	logger := getDBLogger() // Get DB logger

	if &logger == new(loggo.Logger) { // Check logger is nil
		t.Fatal("logger is nil") // Panic
	}
}

/* END INTERNAL METHODS TESTS */
