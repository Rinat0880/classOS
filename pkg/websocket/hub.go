package websocket

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients by channel
	// Key: channel name (e.g., "agent::uuid" or "admin::dashboard")
	channels map[string]map[*Client]bool

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast message to specific channel
	broadcast chan *BroadcastMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// BroadcastMessage represents a message to broadcast to a specific channel
type BroadcastMessage struct {
	Channel string
	Message []byte
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		channels:   make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	logrus.Info("WebSocket Hub started")

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToChannel(message)
		}
	}
}

// registerClient registers a new client to a channel
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.channels[client.channel] == nil {
		h.channels[client.channel] = make(map[*Client]bool)
	}
	h.channels[client.channel][client] = true

	logrus.WithFields(logrus.Fields{
		"channel":   client.channel,
		"client_id": client.ID,
		"user_type": client.UserType,
	}).Info("Client registered to channel")
}

// unregisterClient unregisters a client from its channel
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.channels[client.channel]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			// Remove channel if empty
			if len(clients) == 0 {
				delete(h.channels, client.channel)
			}

			logrus.WithFields(logrus.Fields{
				"channel":   client.channel,
				"client_id": client.ID,
				"user_type": client.UserType,
			}).Info("Client unregistered from channel")
		}
	}
}

// broadcastToChannel sends a message to all clients in a specific channel
func (h *Hub) broadcastToChannel(msg *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.channels[msg.Channel]
	if !ok {
		logrus.WithField("channel", msg.Channel).Warn("Channel not found for broadcast")
		return
	}

	for client := range clients {
		select {
		case client.send <- msg.Message:
		default:
			// Client's send buffer is full, close and unregister
			go func(c *Client) {
				h.unregister <- c
			}(client)
		}
	}

	logrus.WithFields(logrus.Fields{
		"channel":      msg.Channel,
		"client_count": len(clients),
	}).Debug("Message broadcast to channel")
}

// BroadcastToChannel broadcasts a message to a specific channel
func (h *Hub) BroadcastToChannel(channel string, message []byte) {
	h.broadcast <- &BroadcastMessage{
		Channel: channel,
		Message: message,
	}
}

// GetChannelClients returns the number of clients in a channel
func (h *Hub) GetChannelClients(channel string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.channels[channel]; ok {
		return len(clients)
	}
	return 0
}

// GetAllChannels returns a list of all active channels
func (h *Hub) GetAllChannels() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	channels := make([]string, 0, len(h.channels))
	for channel := range h.channels {
		channels = append(channels, channel)
	}
	return channels
}

// CleanupStaleConnections removes connections that haven't sent heartbeat
func (h *Hub) CleanupStaleConnections(timeout time.Duration) {
	ticker := time.NewTicker(timeout / 2)
	defer ticker.Stop()

	for range ticker.C {
		h.mu.RLock()
		var staleClients []*Client

		for _, clients := range h.channels {
			for client := range clients {
				if time.Since(client.lastHeartbeat) > timeout {
					staleClients = append(staleClients, client)
				}
			}
		}
		h.mu.RUnlock()

		// Unregister stale clients
		for _, client := range staleClients {
			logrus.WithFields(logrus.Fields{
				"channel":        client.channel,
				"client_id":      client.ID,
				"last_heartbeat": client.lastHeartbeat,
			}).Warn("Removing stale connection")
			h.unregister <- client
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}
