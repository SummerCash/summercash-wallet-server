// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"encoding/json"
	"fmt"

	"github.com/SummerCash/go-summercash/types"
	"github.com/SummerCash/summercash-wallet-server/common"

	"github.com/SummerCash/summercash-wallet-server/accounts"

	"github.com/valyala/fasthttp"
)

// calcBalanceResponse represents a response to a CalcBalance request.
type calcBalanceResponse struct {
	Balance float64 `json:"balance"` // Account balance
}

// getUserTransactionsResponse represents a response to a GetUserTransactions request.
type getUserTransactionsResponse struct {
	Transactions []*types.Transaction `json:"transactions"` // Account transactions
}

// authenticateUserResponse represents a response to an AuthenticateUser request.
type authenticateUserResponse struct {
	Authenticated bool `json:"authenticated"` // Authenticated
}

/* BEGIN EXPORTED METHODS */

// SetupAccountRoutes sets up all account api-related routes.
func (api *JSONHTTPAPI) SetupAccountRoutes() error {
	accountsAPIRoot := "/api/accounts" // Get accounts API root path

	api.Router.POST(fmt.Sprintf("%s/:username", accountsAPIRoot), api.NewAccount)                      // Set NewAccount post
	api.Router.PUT(fmt.Sprintf("%s/:username", accountsAPIRoot), api.RestAccountPassword)              // Set ResetAccountPassword put
	api.Router.GET(fmt.Sprintf("%s/:username", accountsAPIRoot), api.QueryAccount)                     // Set QueryAccount get
	api.Router.GET(fmt.Sprintf("%s/:username/balance", accountsAPIRoot), api.CalculateAccountBalance)  // Set CalculateAccountBalance get
	api.Router.GET(fmt.Sprintf("%s/:username/transactions", accountsAPIRoot), api.GetUserTransactions) // Set GetUserTransactions get
	api.Router.POST(fmt.Sprintf("%s/:username/authenticate", accountsAPIRoot), api.AuthenticateUser)   // Set AuthenticateUser post

	return nil // No error occurred, return nil
}

// NewAccount handles a NewAccount request.
func (api *JSONHTTPAPI) NewAccount(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Enable CORS

	var account *accounts.Account // Initialize account buffer
	var err error                 // Initialize error buffer

	if address := common.GetCtxValue(ctx, "address"); address != nil { // Check address specified
		account, err = api.AccountsDatabase.AddNewAccount(ctx.UserValue("username").(string), string(common.GetCtxValue(ctx, "password")), string(address)) // Add user
	} else {
		account, err = api.AccountsDatabase.CreateNewAccount(ctx.UserValue("username").(string), string(common.GetCtxValue(ctx, "password"))) // Create new account
	}

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewAccount request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, account.String()) // Respond with account string
}

// RestAccountPassword handles a ResetAccountPassword request.
func (api *JSONHTTPAPI) RestAccountPassword(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Enable CORS

	err := api.AccountsDatabase.ResetAccountPassword(ctx.UserValue("username").(string), string(common.GetCtxValue(ctx, "old_password")), string(common.GetCtxValue(ctx, "new_password"))) // Reset password

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
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Enable CORS

	account, err := api.AccountsDatabase.QueryAccountByUsername(ctx.UserValue("username").(string)) // Query account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling QueryAccount request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, account.String()) // Respond with account string
}

// CalculateAccountBalance handles a CalculateAccountBalance request.
func (api *JSONHTTPAPI) CalculateAccountBalance(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Enable CORS

	balance, err := api.AccountsDatabase.GetUserBalance(ctx.UserValue("username").(string)) // Get balance

	if err != nil { // Check for errors
		logger.Errorf("errored while handling GetUserBalance request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	balanceResponse := &calcBalanceResponse{
		Balance: balance, // Set balance
	} // Initialize balance response

	fmt.Fprintf(ctx, balanceResponse.string()) // Respond with balance response instance
}

// GetUserTransactions handles a GetUserTransactions request.
func (api *JSONHTTPAPI) GetUserTransactions(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Enable CORS

	userTransactions, err := api.AccountsDatabase.GetUserTransactions(ctx.UserValue("username").(string)) // Get user transactions

	if err != nil { // Check for errors
		logger.Errorf("errored while handling GetUserTransactions request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // panic
	}

	getUserTransactionsResponse := &getUserTransactionsResponse{
		Transactions: userTransactions, // Set user transactions
	} // Initialize user txs response

	fmt.Fprintf(ctx, getUserTransactionsResponse.string()) // Respond with user transactions response instance
}

// AuthenticateUser handles an AuthenticateUser request.
func (api *JSONHTTPAPI) AuthenticateUser(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Enable CORS

	authenticateUserResponse := &authenticateUserResponse{
		Authenticated: api.AccountsDatabase.Auth(ctx.UserValue("username").(string), string(common.GetCtxValue(ctx, "password"))), // Set authenticated
	}

	fmt.Fprintf(ctx, authenticateUserResponse.string()) // Respond with user authenticate response instance
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// string marshals a calcBalanceResponse into a JSON-formatted string.
func (response *calcBalanceResponse) string() string {
	marshaledVal, _ := json.MarshalIndent(*response, "", "  ") // Marshal value

	return string(marshaledVal) // Return value
}

// string marshals a getUserTransactionsResponse into a JSON-formatted string.
func (response *getUserTransactionsResponse) string() string {
	marshaledval, _ := json.MarshalIndent(*response, "", "  ") // Marshal value

	return string(marshaledval) // Return value
}

// string marshals an authenticateUserResponse into a JSON-formatted string.
func (response *authenticateUserResponse) string() string {
	marshaledval, _ := json.MarshalIndent(*response, "", "  ") // Marshal value

	return string(marshaledval) // Return value
}

/* END INTERNAL METHODS */
