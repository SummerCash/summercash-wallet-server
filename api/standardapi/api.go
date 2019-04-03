// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"github.com/SummerCash/summercash-wallet-server/accounts"
	"github.com/buaazp/fasthttprouter"
	"github.com/juju/loggo"
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

	AccountsDatabase *accounts.DB `json:"-"` // Accounts database
}

/* BEGIN EXPORTED METHODS */

// NewJSONHTTPAPI initializes a new JSONHTTPAPI instance.
func NewJSONHTTPAPI(baseURI string, provider string, accountsDB *accounts.DB) *JSONHTTPAPI {
	return &JSONHTTPAPI{
		BaseURI:          baseURI,    // Set base URI
		Provider:         provider,   // Set provider
		AccountsDatabase: accountsDB, // Set accounts DB
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
	api.Router = fasthttprouter.New() // Initialize router

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// getAPILogger gets the API package logger.
func getAPILogger() loggo.Logger {
	logger := loggo.GetLogger("API") // Get logger

	loggo.ConfigureLoggers("API=INFO") // Configure loggers

	return logger // Return logger
}

/* END INTERNAL METHODS */
