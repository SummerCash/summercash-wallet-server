// Package api defines the summercash-wallet-server API.
package api

// API defines a summercash-wallet-server API.
type API interface {
	GetBaseURI() string // GetBaseURI gets the APIs base URI.

	GetAvailableAPIs() []string // GetAvailableAPIs gets all available APIs.

	GetServingProtocol() string // GetServingProtocol returns the APIs serving protocol (e.g. https/2).

	GetInputType() string // GetInputType returns the APIs input type (e.g. form data, JSON, protobuf).

	GetResponseType() string // GetResponseType returns the APIs response type (e.g. JSON).

	StartServing() error // Start serving starts serving the given API.
}
