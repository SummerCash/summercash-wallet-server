// Package main is the summercash-wallet-server entry point.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SummerCash/summercash-wallet-server/accounts"
	"github.com/SummerCash/summercash-wallet-server/api/standardapi"

	"github.com/SummerCash/summercash-wallet-server/common"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"

	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/p2p"
	"github.com/SummerCash/go-summercash/rpc"
	"github.com/SummerCash/go-summercash/validator"
)

var (
	nodeRPCPortFlag = flag.Int("node-rpc-port", 8080, "starts the go-summercash RPC server on a given port") // Init node rpc port flag
	nodePortFlag    = flag.Int("node-port", 3000, "starts the go-summercash node on a given port")           // Init node port flag
	networkFlag     = flag.String("network", "main_net", "starts the go-summercash node on a given network") // Init network flag

	logger = loggo.GetLogger("") // Get logger

	logFile *os.File // Log file
)

// main is the summercash-wallet-server entry point.
func main() {
	flag.Parse() // Parse flags

	err := configureLogger() // Configure logger

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	err = startSummercashRPCServer() // Start summercash RPC server

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	err = startNode() // Start node

	if err != nil { // Check for errors
		logger.Criticalf("main panicked: %s", err.Error()) // Log pending panic

		panic(err) // Panic
	}
}

// configureLogger configures the summercash-wallet-server logger.
func configureLogger() error {
	err := common.CreateDirIfDoesNotExit(filepath.FromSlash(common.LogsDir)) // Create log dir

	if err != nil { // Check for errors
		return err // Return found error
	}

	logFile, err = os.OpenFile(filepath.FromSlash(fmt.Sprintf("%s/logs_%s.txt", common.LogsDir, time.Now().Format("2006-01-02_15-04-05"))), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // Create log file

	if err != nil { // Check for errors
		return err // Return found error
	}

	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr)) // Enabled colored output

	loggo.RegisterWriter("logs", loggo.NewSimpleWriter(logFile, loggo.DefaultFormatter)) // Register file writer

	return nil // No error occurred, return nil
}

// startSummercashRPCServer starts the go-summercash RPC server.
func startSummercashRPCServer() error {
	return rpc.StartRPCServer(*nodeRPCPortFlag) // Start RPC server
}

// startNode starts a new go-summercash node.
func startNode() error {
	ctx, cancel := context.WithCancel(context.Background()) // Get node context

	defer cancel() // Cancel

	host, err := p2p.NewHost(ctx, *nodePortFlag) // Initialize libp2p host with context and nat manager

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	config, err := config.ReadChainConfigFromMemory() // Read chain config

	if err != nil { // Check for errors
		config, err = p2p.BootstrapConfig(ctx, host, p2p.GetBestBootstrapAddress(ctx, host, *networkFlag), *networkFlag) // Bootstrap config

		if err != nil { // Check for errors
			panic(err) // panic
		}

		err = config.WriteToMemory() // Write config to memory

		if err != nil { // Check for errors
			panic(err) // Panic
		}
	}

	validator := validator.Validator(validator.NewStandardValidator(config)) // Initialize validator

	client := p2p.NewClient(host, &validator, *networkFlag) // Initialize client

	err = client.StartServingStreams() // Start serving

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	if p2p.GetBestBootstrapAddress(ctx, host, *networkFlag) != "localhost" { // Check can sync
		err = client.SyncNetwork() // Sync network

		if err != nil { // Check for errors
			panic(err) // Panic
		}
	}

	go client.StartIntermittentSync(60 * time.Second) // Start intermittent sync

	return nil // No error occurred, return nil
}

// startServingStandardHTTPJSONAPI starts serving the standard HTTP JSON API.
func startServingStandardHTTPJSONAPI() error {
	db, err := accounts.OpenDB() // Open db

	if err != nil { // Check for errors
		return err // Return found error
	}

	api := standardapi.NewJSONHTTPAPI("http://localhost:8000/api", "", db) // Initialize API instance

	err = api.StartServing() // Start serving

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}
