// Package common outlines common helper methods and types.
package common

import (
	"testing"

	"github.com/valyala/fasthttp"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestGetCtxValue tests the functionality of the GetCtxValue() helper method.
func TestGetCtxValue(t *testing.T) {
	nilCtx := &fasthttp.RequestCtx{} // Init nil ctx

	if value := GetCtxValue(nilCtx, "test"); value != nil { // Check for errors
		t.Fatal("should not have been able to obtain value of nil ctx") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
