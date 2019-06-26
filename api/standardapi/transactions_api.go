// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NaySoftware/go-fcm"
	"github.com/valyala/fasthttp"

	summercashCommon "github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/SummerCash/summercash-wallet-server/transactions"
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
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")             // Allow CORS
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header
	ctx.Response.Header.Set("Content-Type", "application/json")             // Set content type

	var recipient summercashCommon.Address // Init recipient buffer
	var err error                          // Init error buffer

	if string(common.GetCtxValue(ctx, "username")) == "faucet" { // Check wants to send from faucet
		logger.Errorf("user with address %s tried to send tx from faucet account", ctx.RemoteAddr().String()) // Log error

		panic(errors.New("cannot send transaction from faucet wallet")) // Panic
	}

	if !strings.Contains(string(common.GetCtxValue(ctx, "recipient")), "0x") { // Check is sending to username
		recipientAccount, err := api.AccountsDatabase.QueryAccountByUsername(string(common.GetCtxValue(ctx, "recipient"))) // Query account

		if err != nil { // Check for errors
			logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

			panic(err) // Panic
		}

		recipient = recipientAccount.Address // Set address
	} else {
		recipient, err = summercashCommon.StringToAddress(string(common.GetCtxValue(ctx, "recipient"))) // Parse recipient

		if err != nil { // Check for errors
			logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

			panic(err) // Panic
		}
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

	if !strings.Contains(string(common.GetCtxValue(ctx, "recipient")), "0x") && os.Getenv("FCM_KEY") != "" { // Check is username recipient
		recipientAccount, err := api.AccountsDatabase.QueryAccountByUsername(string(common.GetCtxValue(ctx, "recipient"))) // Query account

		if err != nil { // Check for errors
			logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

			panic(err) // Panic
		}

		amount, _ := transaction.Amount.Float64() // Get tx amount

		data := map[string]string{
			"msg": "New Transaction",
			"sum": fmt.Sprintf("Received %f SMC from %s.", amount, transaction.Sender.String()),
		}

		if api.WebsocketManager != nil && api.UseWebsocket { // Check uses websockets
			recipient := string(common.GetCtxValue(ctx, "recipient")) // Get recipient username

			sender := string(common.GetCtxValue(ctx, "username")) // Get sender username

			if !strings.Contains(recipient, "0x") { // Check recipient has username
				recipientBalance, err := api.AccountsDatabase.GetUserBalance(recipient) // Calculate recipient balance

				if err != nil { // Check for errors
					logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

					panic(err) // Panic
				}

				recipientFloatBalance, _ := recipientBalance.Float64() // Get float value

				payload := []byte(fmt.Sprintf("%f:%s", recipientFloatBalance, transaction.String())) // Initialize payload

				for _, session := range api.WebsocketManager.Clients[recipient] { // Iterate through recipient WS sessions
					session.Write(payload) // Write payload
				}
			}

			if !strings.Contains(sender, "0x") { // Check sender has username
				senderBalance, err := api.AccountsDatabase.GetUserBalance(sender) // Calculate sender balance

				if err != nil { // Check for errors
					logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

					panic(err) // Panic
				}

				senderFloatBalance, _ := senderBalance.Float64() // Get float value

				payload := []byte(fmt.Sprintf("%f:%s", senderFloatBalance, transaction.String())) // Initialize payload

				for _, session := range api.WebsocketManager.Clients[sender] { // Iterate through sender WS sessions
					session.Write(payload) // Write payload
				}
			}
		}

		client := fcm.NewFcmClient(os.Getenv("FCM_KEY")) // Init client

		client.NewFcmRegIdsMsg(recipientAccount.FcmTokens, data) // Init message

		_, err = client.Send() // Send notification

		if err != nil { // Check for errors
			logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error
		}
	}

	fmt.Fprintf(ctx, transaction.String()) // Write tx string value
}

/* END EXPORTED METHODS */
