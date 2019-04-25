// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
	Transactions []*types.StringTransaction `json:"transactions"` // Account transactions
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
	api.Router.DELETE(fmt.Sprintf("%s/:username", accountsAPIRoot), api.DeleteUser)                    // Set DeleteUser delete

	return nil // No error occurred, return nil
}

// NewAccount handles a NewAccount request.
func (api *JSONHTTPAPI) NewAccount(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

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
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

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
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

	account, err := api.AccountsDatabase.QueryAccountByUsername(ctx.UserValue("username").(string)) // Query account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling QueryAccount request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, account.String()) // Respond with account string
}

// CalculateAccountBalance handles a CalculateAccountBalance request.
func (api *JSONHTTPAPI) CalculateAccountBalance(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

	balance, err := api.AccountsDatabase.GetUserBalance(ctx.UserValue("username").(string)) // Get balance

	if err != nil { // Check for errors
		logger.Errorf("errored while handling GetUserBalance request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	floatVal, _ := balance.Float64() // Get float val

	balanceResponse := &calcBalanceResponse{
		Balance: floatVal, // Set balance
	} // Initialize balance response

	fmt.Fprintf(ctx, balanceResponse.string()) // Respond with balance response instance
}

// GetUserTransactions handles a GetUserTransactions request.
func (api *JSONHTTPAPI) GetUserTransactions(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

	userTransactions, err := api.AccountsDatabase.GetUserTransactions(ctx.UserValue("username").(string)) // Get user transactions

	if err != nil { // Check for errors
		logger.Errorf("errored while handling GetUserTransactions request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // panic
	}

	var stringTransactions []*types.StringTransaction // Init string tx buffer

	for _, transaction := range userTransactions { // Iterate through user txs
		stringTransaction := &types.StringTransaction{
			AccountNonce:            transaction.AccountNonce,                                                   // Set account nonce
			SenderHex:               transaction.Sender.String(),                                                // Set sender hex
			RecipientHex:            transaction.Recipient.String(),                                             // Set recipient hex
			Amount:                  transaction.Amount.String(),                                                // Set amount
			Payload:                 transaction.Payload,                                                        // Set payload
			Signature:               transaction.Signature,                                                      // Set signature
			ParentTx:                transaction.ParentTx,                                                       // Set parent
			Timestamp:               transaction.Timestamp.Add(-4 * time.Hour).Format("01/02/2006 03:04:05 PM"), // Set timestamp
			DeployedContractAddress: transaction.DeployedContractAddress,                                        // Set deployed contract address
			ContractCreation:        transaction.ContractCreation,                                               // Set is contract creation
			Genesis:                 transaction.Genesis,                                                        // Set is genesis
			Logs:                    transaction.Logs,                                                           // Set logs
			HashHex:                 transaction.Hash.String(),                                                  // Set hash hex
		} // Init string transaction

		stringTransactions = append(stringTransactions, stringTransaction) // Append string tx
	}

	getUserTransactionsResponse := &getUserTransactionsResponse{
		Transactions: stringTransactions, // Set string transactions
	} // Initialize user txs response

	fmt.Fprintf(ctx, getUserTransactionsResponse.string()) // Respond with user transactions response instance
}

// AuthenticateUser handles an AuthenticateUser request.
func (api *JSONHTTPAPI) AuthenticateUser(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

	if !api.AccountsDatabase.Auth(ctx.UserValue("username").(string), string(common.GetCtxValue(ctx, "password"))) { // Check cannot authenticate
		logger.Errorf("errored while handling AuthenticateUser request with username %s", ctx.UserValue("username")) // Log error

		panic(errors.New("invalid username or password")) // panic
	}

	account, err := api.AccountsDatabase.QueryAccountByUsername(ctx.UserValue("username").(string)) // Get account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling AuthenticateUser request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // panic
	}

	fmt.Fprintf(ctx, account.String()) // Respond with user details
}

// DeleteUser handles a DeleteUser request.
func (api *JSONHTTPAPI) DeleteUser(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*") // Allow CORS

	err := api.AccountsDatabase.DeleteAccount(ctx.UserValue("username").(string), string(common.GetCtxValue(ctx, "password"))) // Delete account

	if err != nil { // Check for errors
		logger.Errorf("errored while handling DeleteUser request with username %s: %s", ctx.UserValue("username"), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, fmt.Sprintf("{%smessage%s: %sAccount deleted successfully%s}", `"`, `"`, `"`, `"`)) // Respond with success
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
