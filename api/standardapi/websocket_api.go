// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

// ConnectionManager manages a set of active websocket connections.
type ConnectionManager struct {
	Clients map[string][]*melody.Session `json:"clients"` // Connected clients
}

/* BEGIN EXPORTED METHODS */

// SetupWebsocketRoutes sets up all the websocket api-related routes.
func (api *JSONHTTPAPI) SetupWebsocketRoutes() error {
	websocketAPIRoot := "/ws/:username" // Get websocket API root path.

	api.WebsocketManager = &ConnectionManager{
		Clients: make(map[string][]*melody.Session), // Set client manager
	}

	api.MiscAPIRouter.GET(websocketAPIRoot, api.HandleWebsocketGet) // Set /ws handler

	api.Melody.HandleConnect(api.HandleConnection) // Set new WebSocket conn handler

	return nil // No error occurred, return nil
}

// HandleWebsocketGet handles an incoming GET request for the /ws path.
func (api *JSONHTTPAPI) HandleWebsocketGet(c *gin.Context) {
	logger.Infof("handling WebSocket GET request from addr %s", c.ClientIP()) // Log accepted request

	err := api.Melody.HandleRequest(c.Writer, c.Request) // Handle request

	if err != nil { // Check for errors
		logger.Errorf("errored while handling WebSocket GET request from addr %s: %s", c.ClientIP(), err.Error()) // Log found error
	}
}

// HandleConnection handles an incoming WebSocket connection.
func (api *JSONHTTPAPI) HandleConnection(s *melody.Session) {
	splitURL := strings.Split(s.Request.URL.String(), "/") // Split URL

	username := splitURL[len(splitURL)-1] // Get last element in split URL

	api.WebsocketManager.Clients[username] = append(api.WebsocketManager.Clients[username], s) // Append session to user sessions
}

// HandleDisconnect handles a disconnected WebSocket connection.
func (api *JSONHTTPAPI) HandleDisconnect(s *melody.Session) {
	splitURL := strings.Split(s.Request.URL.String(), "/") // Split URL

	username := splitURL[len(splitURL)-1] // Get last element in split URL

	for i, currentSession := range api.WebsocketManager.Clients[username] { // Iterate through listening clients in scope of username
		if currentSession == s { // Check is same session
			api.WebsocketManager.Clients[username][i] = api.WebsocketManager.Clients[username][len(api.WebsocketManager.Clients[username])-1] // Clear slot
			api.WebsocketManager.Clients[username] = api.WebsocketManager.Clients[username][:len(api.WebsocketManager.Clients[username])-1]   // Remove from list of listening clients
		}
	}
}

/* END EXPORTED METHODS */
