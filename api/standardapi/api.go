// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"github.com/SummerCash/summercash-wallet-server/accounts"
	"github.com/buaazp/fasthttprouter"
)

// StandardJSONHTTPAPI is an instance of an API providing the standard API set via https/2 JSON.
type StandardJSONHTTPAPI struct {
	BaseURI string `json:"base_uri"` // Base URI

	Provider string `json:"provider"` // Node provider

	Router *fasthttprouter.Router `json:"-"` // Router

	AccountsDatabase *accounts.DB `json:"-"` // Accounts database
}

/* BEGIN EXPORTED METHODS */

// NewStandardJSONHTTPAPI initializes a new StandardJSONHTTPAPI instance.
func NewStandardJSONHTTPAPI(baseURI string, provider string, accountsDB *accounts.DB) *StandardJSONHTTPAPI {
	return &StandardJSONHTTPAPI{
		BaseURI:          baseURI,    // Set base URI
		Provider:         provider,   // Set provider
		AccountsDatabase: accountsDB, // Set accounts DB
	}
}

// GetAvailableAPIs gets the available APIs.
func (api *StandardJSONHTTPAPI) GetAvailableAPIs() []string {
	return []string{} // Return available APIs
}

// GetServingProtocol gets the serving protocol.
func (api *StandardJSONHTTPAPI) GetServingProtocol() string {
	return "https/2" // Return protocol
}

// GetInputType gets the input type.
func (api *StandardJSONHTTPAPI) GetInputType() string {
	return "JSON" // Return input type
}

// GetResponseType gets the response type.
func (api *StandardJSONHTTPAPI) GetResponseType() string {
	return "JSON" // Return response type
}

// StartServing starts serving the API.
func (api *StandardJSONHTTPAPI) StartServing() error {
	api.Router = fasthttprouter.New() // Initialize router

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
