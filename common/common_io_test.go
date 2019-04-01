// Package common outlines common helper methods and types.
package common

import "testing"

/* BEGIN EXPORTED METHODS TESTS */

// TestCreateDirIfDoesNotExist tests the CreateDirIfDoesNotExit() helper method.
func TestCreateDirIfDoesNotExit(t *testing.T) {
	err := CreateDirIfDoesNotExit("test") // Create dir

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
