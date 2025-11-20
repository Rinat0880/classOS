package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "github.com/rinat0880/classOS_backend/pkg/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Add proper origin checking in production
		return true
	},
}

// serveWs handles websocket requests from the peer.
func (h *Handler) serveWs(c *gin.Context) {
	// Get client info from context (set by auth middleware)
	clientID, exists := c.Get("ws_client_id")
	if !exists {
		logrus.Error("ws_client_id not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userType, exists := c.Get("ws_user_type")
	if !exists {
		logrus.Error("ws_user_type not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	channel, exists := c.Get("ws_channel")
	if !exists {
		logrus.Error("ws_channel not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade connection to WebSocket")
		return
	}

	// Create new client
	client := ws.NewClient(
		clientID.(string),
		userType.(string),
		channel.(string),
		h.wsHub,
		conn,
	)

	// Register client with hub
	h.wsHub.RegisterClient(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()

	logrus.WithFields(logrus.Fields{
		"client_id": clientID,
		"user_type": userType,
		"channel":   channel,
	}).Info("WebSocket connection established")
}

// getWsStatus returns WebSocket server status
func (h *Handler) getWsStatus(c *gin.Context) {
	channels := h.wsHub.GetAllChannels()

	channelStats := make(map[string]int)
	for _, ch := range channels {
		channelStats[ch] = h.wsHub.GetChannelClients(ch)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "running",
		"total_channels": len(channels),
		"channels":       channelStats,
	})
}
