// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

/* BEGIN EXPORTED METHODS */

// SetupAccountRoutes sets up all account api-related routes.
func (api *JSONHTTPAPI) SetupAccountRoutes() error {
	accountsAPIRoot := fmt.Sprintf("%s/accounts", api.BaseURI) // Get accounts API root path

	api.Router.POST(fmt.Sprintf("%s/:username", accountsAPIRoot), api.NewAccount) // Set NewAccount post

	return nil // No error occurred, return nil
}

// NewAccount handles a NewAccount request.
func (api *JSONHTTPAPI) NewAccount(ctx *fasthttp.RequestCtx) {
	account, err := api.AccountsDatabase.CreateNewAccount(ctx.UserValue("username").(string), string(ctx.FormValue("password"))) // Create new account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewAccount request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		return // Return
	}

	fmt.Fprintf(ctx, account.String()) // Respond with account string
}

/* END EXPORTED METHODS */
