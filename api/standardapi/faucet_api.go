// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"
	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/valyala/fasthttp"
	"math/big"
	"time"
)

/* BEGIN EXPORTED METHODS */

// SetupFaucetRoutes sets up all the faucet api-related routes.
func (api *JSONHTTPAPI) SetupFaucetRoutes() error {
	faucetAPIRoot := "/api/faucet" // Get faucet API root path

	api.Router.POST(fmt.Sprintf("%s/Claim", faucetAPIRoot), api.Claim)                  // Set Claim post
	api.Router.GET(fmt.Sprintf("%s/:username/NextClaim", faucetAPIRoot), api.NextClaim) // Set Claim post

	return nil // No error occurred, return nil
}

// Claim handles a Claim request.
func (api *JSONHTTPAPI) Claim(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

	account, err := api.AccountsDatabase.QueryAccountByUsername(string(common.GetCtxValue(ctx, "username"))) // Query account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling Claim request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	amount, _, _ := big.ParseFloat(string(common.GetCtxValue(ctx, "amount")), 10, 350, big.ToNearestEven) // Parse amount

	if err != nil { // Check for errors
		logger.Errorf("errored while handling Claim request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	err = (*api.Faucet).Claim(account, amount) // Claim

	if err != nil { // Check for errors
		logger.Errorf("errored while handling Claim request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, fmt.Sprintf("{%smessage%s: %sFaucet bounty claimed successfully%s}", `"`, `"`, `"`, `"`)) // Write response
}

// NextClaim handles a NextClaim request.
func (api *JSONHTTPAPI) NextClaim(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

	account, err := api.AccountsDatabase.QueryAccountByUsername(string(common.GetCtxValue(ctx, "username"))) // Query account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling Claim request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	var timeUntilNextClaim = time.Now().Sub(time.Now()) // Init buffer

	if !(*api.Faucet).AccountCanClaim(account) { // Check can claim
		timeUntilNextClaim = (*api.Faucet).AccountLastClaim(account).Sub(time.Now()) // Get time until next claim
	}

	fmt.Fprintf(ctx, fmt.Sprintf("{%stime%s: %s%s%s}", `"`, `"`, `"`, fmt.Sprintf("%d:%d:%d", uint(timeUntilNextClaim.Hours()), uint(timeUntilNextClaim.Minutes()), uint(timeUntilNextClaim.Seconds())), `"`)) // Write time until
}

/* END EXPORTED METHODS */
