// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	fasthttprouter "github.com/fasthttp/router"
	"github.com/juju/loggo"
	"github.com/valyala/fasthttp"

	"github.com/SummerCash/summercash-wallet-server/accounts"
	"github.com/SummerCash/summercash-wallet-server/faucet"
)

var (
	// logger is the api package logger.
	logger = getAPILogger()
)

// JSONHTTPAPI is an instance of an API providing the standard API set via https/2 JSON.
type JSONHTTPAPI struct {
	BaseURI string `json:"base_uri"` // Base URI

	Provider string `json:"provider"` // Node provider

	Router *fasthttprouter.Router `json:"-"` // Router

	ContentRouter *fasthttprouter.Router `json:"-"` // Content router

	AccountsDatabase *accounts.DB `json:"-"` // Accounts database

	Faucet *faucet.Faucet `json:"-"` // Faucet

	ContentDir string `json:"content_dir"` // Static content directory
}

// errorResponse represents a JSON error.
type errorResponse struct {
	Error string `json:"error"` // Error
}

/* BEGIN EXPORTED METHODS */

// NewJSONHTTPAPI initializes a new JSONHTTPAPI instance.
func NewJSONHTTPAPI(baseURI string, provider string, accountsDB *accounts.DB, faucet *faucet.Faucet, contentDir string) *JSONHTTPAPI {
	return &JSONHTTPAPI{
		BaseURI:          baseURI,    // Set base URI
		Provider:         provider,   // Set provider
		AccountsDatabase: accountsDB, // Set accounts DB
		ContentDir:       contentDir, // Set content dir
		Faucet:           faucet,     // Set faucet
	}
}

// GetAvailableAPIs gets the available APIs.
func (api *JSONHTTPAPI) GetAvailableAPIs() []string {
	return []string{} // Return available APIs
}

// GetServingProtocol gets the serving protocol.
func (api *JSONHTTPAPI) GetServingProtocol() string {
	return "https/2" // Return protocol
}

// GetInputType gets the input type.
func (api *JSONHTTPAPI) GetInputType() string {
	return "JSON" // Return input type
}

// GetResponseType gets the response type.
func (api *JSONHTTPAPI) GetResponseType() string {
	return "JSON" // Return response type
}

// StartServing starts serving the API.
func (api *JSONHTTPAPI) StartServing() error {
	api.Router = fasthttprouter.New()        // Initialize router
	api.ContentRouter = fasthttprouter.New() // Initialize content router

	var err error // Init error buffer

	serveContent := false // Should serve content

	if api.ContentDir != "" { // Check should serve content
		if _, err := os.Stat(api.ContentDir); !os.IsNotExist(err) { // Check can serve content
			serveContent = true // Set can serve
		}

		api.ContentDir, _ = filepath.Abs(api.ContentDir) // Get absolute path

		api.ContentRouter.ServeFiles("/wallet/*filepath", api.ContentDir) // Serve files
	}

	api.Router.PanicHandler = api.HandlePanic // Set panic handler

	err = api.SetupAccountRoutes() // Start serving accounts API

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = api.SetupTransactionsRoutes() // Start serving transactions API

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = api.SetupFaucetRoutes() // Start serving faucet API

	if err != nil { // Check for errors
		return err // Return found error
	}

	switch serveContent {
	case true:
		go fasthttp.ListenAndServeTLS(strings.Split(api.BaseURI, "/api")[0], "generalCert.pem", "generalKey.pem", api.Router.Handler) // Start serving

		err = fasthttp.ListenAndServeTLS(":443", "generalCert.pem", "generalKey.pem", api.ContentRouter.Handler) // Start serving

		if err != nil { // Check for errors
			return err // Return found error
		}
	default:
		err = fasthttp.ListenAndServeTLS(strings.Split(api.BaseURI, "/api")[0], "generalCert.pem", "generalKey.pem", api.Router.Handler) // Start serving

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	if serveContent { // Check can serve content dir

	}

	return nil // No error occurred, return nil
}

// HandlePanic handles a panic.
func (api *JSONHTTPAPI) HandlePanic(ctx *fasthttp.RequestCtx, panic interface{}) {
	errorInstance := &errorResponse{
		Error: panic.(error).Error(), // Set error
	}

	fmt.Fprintf(ctx, errorInstance.string()) // Log error
	fmt.Println(panic.(error).Error())       // Log error
	debug.PrintStack()                       // Print stack trace
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// string marshals an error response into a string.
func (response *errorResponse) string() string {
	marshaledVal, _ := json.MarshalIndent(*response, "", "  ") // marshal

	return string(marshaledVal) // Return response
}

// getAPILogger gets the API package logger.
func getAPILogger() loggo.Logger {
	logger := loggo.GetLogger("API") // Get logger

	loggo.ConfigureLoggers("API=INFO") // Configure loggers

	return logger // Return logger
}

/* END INTERNAL METHODS */
