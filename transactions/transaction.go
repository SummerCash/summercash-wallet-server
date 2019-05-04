// Package transactions outlines helper methods for the SummerCash tx api.
package transactions

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/types"
	"github.com/SummerCash/go-summercash/validator"
	"github.com/SummerCash/summercash-wallet-server/accounts"

	summercashAccounts "github.com/SummerCash/go-summercash/accounts"

	transactionProto "github.com/SummerCash/go-summercash/intrnl/rpc/proto/transaction"
	transactionServer "github.com/SummerCash/go-summercash/intrnl/rpc/transaction"
)

/* BEGIN EXPORTED METHODS */

// NewTransaction creates, signs, and publishes a new transaction from a given user to a given address.
func NewTransaction(accountsDB *accounts.DB, username string, password string, recipientAddress *common.Address, amount float64, payload []byte) (*types.Transaction, error) {
	account, err := accountsDB.QueryAccountByUsername(username) // Query account

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	if authenticated := accountsDB.Auth(username, password); !authenticated { // Check could not authenticate
		return &types.Transaction{}, errors.New("invalid username or password") // Return found error
	}

	accountChain, err := types.ReadChainFromMemory(account.Address) // Read chain

	if err != nil { // Check for errors
		accountChain, err = types.NewChain(account.Address) // Initialize chain

		if err != nil { // Check for errors
			return &types.Transaction{}, err // Return found error
		}
	}

	fmt.Println("test")

	var parentTransaction *types.Transaction // Init parent tx buffer

	targetNonce := uint64(0.0) // Init target nonce

	if len(accountChain.Transactions) > 0 { // Check has txs
		parentTransaction = accountChain.Transactions[len(accountChain.Transactions)-1] // Set parent transaction

		targetNonce = accountChain.CalculateTargetNonce() // Set nonce
	}

	fmt.Println("test")

	transaction, err := types.NewTransaction(targetNonce, parentTransaction, &account.Address, recipientAddress, big.NewFloat(amount), payload) // Initialize transaction

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	summercashAccount, err := summercashAccounts.ReadAccountFromMemory(account.Address) // Read account

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	fmt.Println("test")

	config, err := config.ReadChainConfigFromMemory() // Read config from memory

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	fmt.Println("test")

	validator := validator.Validator(validator.NewStandardValidator(config)) // Initialize validator

	err = types.SignTransaction(transaction, summercashAccount.PrivateKey) // Sign transaction

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	fmt.Println("test")

	err = validator.ValidateTransaction(transaction) // Validate transaction

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	fmt.Println("test")

	err = transaction.WriteToMemory() // Write tx to memory

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	fmt.Println("test")

	rpcServer := &transactionServer.Server{} // Initialize mock RPC server

	publishCtx, cancel := context.WithCancel(context.Background()) // Get ctx

	defer cancel() // Cancel

	_, err = rpcServer.Publish(publishCtx, &transactionProto.GeneralRequest{Address: transaction.Hash.String()}) // Publish

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	return transaction, nil // Return tx
}

/* END EXPORTED METHODS */
