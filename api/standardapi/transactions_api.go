// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"
	"strconv"

	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/SummerCash/summercash-wallet-server/transactions"

	"github.com/valyala/fasthttp"

	summercashCommon "github.com/SummerCash/go-summercash/common"
)

/* BEGIN EXPORTED METHODS */

// SetupTransactionsRoutes sets up all the transactions api-related routes.
func (api *JSONHTTPAPI) SetupTransactionsRoutes() error {
	transactionsAPIRoot := "/api/transactions" // Get transactions API root path

	api.Router.POST(fmt.Sprintf("%s/NewTransaction", transactionsAPIRoot), api.NewTransaction) // Set NewTransaction post

	return nil // No error occurred, return nil
}

// NewTransaction handles a NewTransaction request.
func (api *JSONHTTPAPI) NewTransaction(ctx *fasthttp.RequestCtx) {
	recipient, err := summercashCommon.StringToAddress(string(common.GetCtxValue(ctx, "recipient"))) // Parse recipient

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	amount, err := strconv.ParseFloat(string(common.GetCtxValue(ctx, "amount")), 64) // Parse amount

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	transaction, err := transactions.NewTransaction(api.AccountsDatabase, string(common.GetCtxValue(ctx, "username")), string(common.GetCtxValue(ctx, "password")), &recipient, amount, common.GetCtxValue(ctx, "payload")) // Initialize transaction

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, transaction.String()) // Write tx string value
}

/* END EXPORTED METHODS */
