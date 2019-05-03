// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/valyala/fasthttp"
)

/* BEGIN EXPORTED METHODS */

// SetupFaucetRoutes sets up all the faucet api-related routes.
func (api *JSONHTTPAPI) SetupFaucetRoutes() error {
	faucetAPIRoot := "/api/faucet" // Get faucet API root path

	api.Router.POST(fmt.Sprintf("%s/Claim", faucetAPIRoot), api.Claim)                              // Set Claim post
	api.Router.GET(fmt.Sprintf("%s/:username/NextClaimTime", faucetAPIRoot), api.NextClaim)         // Set Claim get
	api.Router.GET(fmt.Sprintf("%s/:username/NextClaimAmount", faucetAPIRoot), api.NextClaimAmount) // Set Claim amount get

	return nil // No error occurred, return nil
}

// Claim handles a Claim request.
func (api *JSONHTTPAPI) Claim(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")             // Allow CORS
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header
	ctx.Response.Header.Set("Content-Type", "application/json")             // Set content type

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
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")             // Allow CORS
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header
	ctx.Response.Header.Set("Content-Type", "application/json")             // Set content type

	account, err := api.AccountsDatabase.QueryAccountByUsername(string(common.GetCtxValue(ctx, "username"))) // Query account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NextClaim request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	timeUntilNextClaim := time.Duration(0) // Init buffer

	var hours, minutes, seconds string // Init time until string buffers

	if !(*api.Faucet).AccountCanClaim(account) { // Check can claim
		timeUntilNextClaim = time.Until((*api.Faucet).AccountLastClaim(account).Add((*(*api.Faucet).GetRuleset()).GetClaimPeriod())) // Get time until next claim

		hours = strings.Split(timeUntilNextClaim.String(), "h")[0]                                                                        // Get hours until
		minutes = strings.Split(strings.Split(timeUntilNextClaim.String(), "h")[1], "m")[0]                                               // Get minutes until
		seconds = strings.Split(strings.Split(strings.Split(strings.Split(timeUntilNextClaim.String(), "h")[1], "m")[1], "s")[0], ".")[0] // Get seconds until
	} else {
		hours, minutes, seconds = "00", "00", "00" // Set to 0
	}

	fmt.Fprintf(ctx, fmt.Sprintf("{%stime%s: %s%s%s}", `"`, `"`, `"`, fmt.Sprintf("%s:%s:%s", hours, minutes, seconds), `"`)) // Write time until
}

// NextClaimAmount handles a NextClaimAmount request.
func (api *JSONHTTPAPI) NextClaimAmount(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")             // Allow CORS
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header
	ctx.Response.Header.Set("Content-Type", "application/json")             // Set content type

	account, err := api.AccountsDatabase.QueryAccountByUsername(string(common.GetCtxValue(ctx, "username"))) // Query account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NextClaimAmount request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	amount := (*api.Faucet).AmountCanClaim(account) // Get amount can claim

	if amount.Cmp(big.NewFloat(0)) == 0 { // Check is zero
		amount = account.LastFaucetClaimAmount // Get last account claim
	}

	intVal, _ := amount.Float64() // Get float value

	fmt.Fprintf(ctx, fmt.Sprintf("{%samount%s: %s%f%s}", `"`, `"`, `"`, intVal, `"`)) // Write time until
}

/* END EXPORTED METHODS */
