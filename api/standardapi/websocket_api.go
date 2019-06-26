// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// ConnectionManager manages a set of active websocket connections.
type ConnectionManager struct {
	Clients map[string]*websocket.Conn `json:"clients"` // Connected clients

	upgrader websocket.Upgrader // Connection ugprader
}

/* BEGIN EXPORTED METHODS */

// SetupWebsocketRoutes sets up all the websocket api-related routes.
func (api *JSONHTTPAPI) SetupWebsocketRoutes() error {
	websocketAPIRoot := "/ws" // Get websocket API root path.

	api.WebsocketManager = &ConnectionManager{
		Clients:  make(map[string]*websocket.Conn), // Set client manager
		upgrader: websocket.Upgrader{},             // Set upgrader
	}

	api.WebsocketManager = new(ConnectionManager) // Set connection manager

	api.Mux.HandleFunc(websocketAPIRoot, api.handleConnection) // Setup websocket route

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// handleConnection handles all incoming websocket connections.
func (api *JSONHTTPAPI) handleConnection(w http.ResponseWriter, r *http.Request) {
	websocket, err := api.WebsocketManager.upgrader.Upgrade(w, r, nil) // Upgrade

	defer websocket.Close() // Close eventually lol

	if err != nil { // Check for errors
		logger.Errorf("errored while handling websocket connection: %s", err.Error()) // Log error

		fmt.Fprintf(w, (&errorResponse{
			Error: err.Error(), // Set error
		}).string()) // Write error to http respwriter

		return // Return
	}

	decoder := json.NewDecoder(r.Body) // Initialize decooder

	jsonMap := make(map[string]*json.RawMessage) // Init JSON map buffer

	err = decoder.Decode(&jsonMap) // Decode

	if err != nil { // Check for errors
		logger.Errorf("errored while handling websocket connection: %s", err.Error()) // Log error

		fmt.Fprintf(w, (&errorResponse{
			Error: err.Error(), // Set error
		}).string()) // Write error to http respwriter

		return // Return
	}

	fmt.Println(jsonMap["username"])
}

/* END INTERNAL METHODS */
