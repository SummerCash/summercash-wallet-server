// Package main is the summercash-wallet-server entry point.
package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"

	summercashCommon "github.com/SummerCash/go-summercash/common"
	"github.com/SummerCash/go-summercash/config"
	"github.com/SummerCash/go-summercash/p2p"
	"github.com/SummerCash/go-summercash/rpc"
	"github.com/SummerCash/go-summercash/validator"
	"github.com/SummerCash/summercash-wallet-server/accounts"
	"github.com/SummerCash/summercash-wallet-server/api/standardapi"
	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/SummerCash/summercash-wallet-server/faucet"
)

var (
	nodeRPCPortFlag   = flag.Int("node-rpc-port", 8080, "starts the go-summercash RPC server on a given port")      // Init node rpc port flag
	nodePortFlag      = flag.Int("node-port", 3000, "starts the go-summercash node on a given port")                // Init node port flag
	networkFlag       = flag.String("network", "main_net", "starts the go-summercash node on a given network")      // Init network flag
	apiPortFlag       = flag.Int("api-port", 2053, "starts api on given port")                                      // Init API port flag
	contentDirFlag    = flag.String("content-dir", filepath.FromSlash("./app"), "serves a given content directory") // Init content dir flag
	dataDirFlag       = flag.String("data-dir", common.DataDir, "starts node with given data directory")            // Init data dir flag
	faucetRewardFlag  = flag.Float64("faucet-reward", 0.00001, "starts faucet api with a given reward amount")      // Init faucet reward flag
	useRemoteNodeFlag = flag.Bool("use-remote-node", false, "skips node start, assumes remote node is up to date")  // Init remote node flag

	logger = loggo.GetLogger("") // Get logger

	logFile *os.File // Log file

	ctx, cancel = context.WithCancel(context.Background()) // Get context
)

// main is the summercash-wallet-server entry point.
func main() {
	flag.Parse() // Parse flags

	common.DataDir = filepath.FromSlash(*dataDirFlag)           // Set data dir
	summercashCommon.DataDir = filepath.FromSlash(*dataDirFlag) // Set data dir

	defer cancel()        // Cancel
	defer logFile.Close() // Close log file

	err := configureLogger() // Configure logger

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	err = startSummercashRPCServer() // Start summercash RPC server

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	if !*useRemoteNodeFlag { // Check must use local node
		err = startNode() // Start node

		if err != nil { // Check for errors
			logger.Criticalf("main panicked: %s", err.Error()) // Log pending panic

			os.Exit(1) // Return
		}
	}

	err = startServingStandardHTTPJSONAPI() // Start serving

	if err != nil { // Check for errors
		logger.Criticalf("main panicked: %s", err.Error()) // Log pending panic

		os.Exit(1) // Return
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
	host, err := p2p.NewHost(ctx, *nodePortFlag, "main_net") // Initialize libp2p host with context and nat manager

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

	c := make(chan os.Signal) // Get control c

	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Notify

	go func() {
		<-c // Wait for ^c

		err = db.DB.Close() // Close dag

		if err != nil { // Check for errors
			logger.Criticalf("db close errored: %s", err.Error()) // Return found error
		}

		logFile.Close() // Close log file

		cancel() // Cancel

		os.Exit(0) // Exit
	}()

	ruleset := faucet.NewStandardRuleset(big.NewFloat(*faucetRewardFlag), 6*time.Hour, []*accounts.Account{}) // Initialize ruleset

	standardFaucet := faucet.NewStandardFaucet(ruleset, db) // Initialize faucet

	abstractFaucet := faucet.Faucet(standardFaucet) // Get interface

	api := standardapi.NewJSONHTTPAPI(fmt.Sprintf(":%d/api", *apiPortFlag), "", db, &abstractFaucet, *contentDirFlag) // Initialize API instance

	err = api.StartServing() // Start serving

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}
