// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"

	"github.com/SummerCash/summercash-wallet-server/accounts"

	"github.com/valyala/fasthttp"
)

/* BEGIN EXPORTED METHODS */

// SetupAccountRoutes sets up all account api-related routes.
func (api *JSONHTTPAPI) SetupAccountRoutes() error {
	accountsAPIRoot := "/api/accounts" // Get accounts API root path

	api.Router.POST(fmt.Sprintf("%s/:username", accountsAPIRoot), api.NewAccount)         // Set NewAccount post
	api.Router.PUT(fmt.Sprintf("%s/:username", accountsAPIRoot), api.RestAccountPassword) // Set ResetAccountPassword put
	api.Router.GET(fmt.Sprintf("%s/:username", accountsAPIRoot), api.QueryAccount)        // Set QueryAccount get

	return nil // No error occurred, return nil
}

// NewAccount handles a NewAccount request.
func (api *JSONHTTPAPI) NewAccount(ctx *fasthttp.RequestCtx) {
	var account *accounts.Account // Initialize account buffer
	var err error                 // Initialize error buffer

	if ctx.FormValue("address") != nil { // Check address specified
		account, err = api.AccountsDatabase.AddNewAccount(ctx.UserValue("username").(string), string(ctx.FormValue("password")), string(ctx.FormValue("address"))) // Add user
	} else {
		account, err = api.AccountsDatabase.CreateNewAccount(ctx.UserValue("username").(string), string(ctx.FormValue("password"))) // Create new account
	}

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewAccount request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, account.String()) // Respond with account string
}

// RestAccountPassword handles a ResetAccountPassword request.
func (api *JSONHTTPAPI) RestAccountPassword(ctx *fasthttp.RequestCtx) {
	err := api.AccountsDatabase.ResetAccountPassword(ctx.UserValue("username").(string), string(ctx.FormValue("old_password")), string(ctx.FormValue("new_password"))) // Reset password

	if err != nil { // Check for errors
		logger.Errorf("errored while handling RestAccountPassword request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	updatedAccount, err := api.AccountsDatabase.QueryAccountByUsername(ctx.UserValue("username").(string)) // Query account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling ResetAccountPassword request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, updatedAccount.String()) // Respond with account string
}

// QueryAccount handles a QueryAccount request.
func (api *JSONHTTPAPI) QueryAccount(ctx *fasthttp.RequestCtx) {
	account, err := api.AccountsDatabase.QueryAccountByUsername(ctx.UserValue("username").(string)) // Query account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling QueryAccount request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, account.String()) // Respond with account string
}

/* END EXPORTED METHODS */
