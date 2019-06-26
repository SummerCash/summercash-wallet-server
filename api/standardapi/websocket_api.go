// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

// ConnectionManager manages a set of active websocket connections.
type ConnectionManager struct {
	Clients map[string]*melody.Session `json:"clients"` // Connected clients
}

/* BEGIN EXPORTED METHODS */

// SetupWebsocketRoutes sets up all the websocket api-related routes.
func (api *JSONHTTPAPI) SetupWebsocketRoutes() error {
	websocketAPIRoot := "/ws/:username" // Get websocket API root path.

	api.WebsocketManager = &ConnectionManager{
		Clients: make(map[string]*melody.Session), // Set client manager
	}

	api.MiscAPIRouter.GET(websocketAPIRoot, api.HandleWebsocketGet) // Set /ws handler

	api.Melody.HandleConnect(api.HandleConnection) // Set new WebSocket conn handler

	return nil // No error occurred, return nil
}

// HandleWebsocketGet handles an incoming GET request for the /ws path.
func (api *JSONHTTPAPI) HandleWebsocketGet(c *gin.Context) {
	api.Melody.HandleRequest(c.Writer, c.Request) // Handle request
}

// HandleConnection handles an incoming WebSocket connection.
func (api *JSONHTTPAPI) HandleConnection(s *melody.Session) {
	fmt.Println(s.Request.URL.Query()["username"][0])
	api.WebsocketManager.Clients[s.Request.URL.Query()["username"][0]] = s // Set session
}

/* END EXPORTED METHODS */
