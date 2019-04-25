// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"
	"github.com/valyala/fasthttp"
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

}

/* END EXPORTED METHODS */
