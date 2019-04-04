// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"

	"github.com/valyala/fasthttp"
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

}

/* END EXPORTED METHODS */
