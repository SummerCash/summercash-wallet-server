package crypto

import (
	"bytes"
	"testing"
)

// TestSalt tests the functionality of the Salt() method.
func TestSalt(t *testing.T) {
	salted := Salt([]byte("test")) // Salt

	if bytes.Equal([]byte("test"), salted) { // Check equal
		t.Fatal("salt should not be same as pass") // Panic
	}
}

// TestVerifySalted tests the functionality of the VerifySalted() method.
func TestVerifySalted(t *testing.T) {
	salted := Salt([]byte("test")) // Salt

	if !VerifySalted(salted, "test") { // Verify salted
		t.Fatal("salt not valid") // Panic
	}
}

// TestSha3 - test functionality of sha3 hashing function
func TestSha3(t *testing.T) {
	hashed := Sha3([]byte("test")) // Hash

	if hashed == nil { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3String - test functionality of sha3 hashing string function
func TestSha3String(t *testing.T) {
	hashed := Sha3String([]byte("test")) // Hash

	if hashed == "" { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3n - test functionality of sha3n hashing function
func TestSha3n(t *testing.T) {
	hashed := Sha3n([]byte("test"), 10) // Hash

	if hashed == nil { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3nString - test functionality of sha3n string hashing function
func TestSha3nString(t *testing.T) {
	hashed := Sha3nString([]byte("test"), 10) // Hash

	if hashed == "" { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3d - test functionality of sha3d hashing function
func TestSha3d(t *testing.T) {
	hashed := Sha3d([]byte("test")) // Hash

	if hashed == nil { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}

// TestSha3dString - test functionality of sha3d hashing string function
func TestSha3dString(t *testing.T) {
	hashed := Sha3dString([]byte("test")) // Hash

	if hashed == "" { // Check is nil
		t.Errorf("invalid hash %s", hashed) // Log found error
		t.FailNow()                         // Panic
	}

	t.Log(hashed) // Log hashed
}
