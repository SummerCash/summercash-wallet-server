// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"
	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/valyala/fasthttp"
	"math/big"
)

/* BEGIN EXPORTED METHODS */

// SetupFaucetRoutes sets up all the faucet api-related routes.
func (api *JSONHTTPAPI) SetupFaucetRoutes() error {
	faucetAPIRoot := "/api/faucet" // Get faucet API root path

	api.Router.POST(fmt.Sprintf("%s/Claim", faucetAPIRoot), api.Claim) // Set Claim post

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
}

/* END EXPORTED METHODS */
