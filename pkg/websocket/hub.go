package websocket

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// BroadcastMessage represents a message to broadcast to a specific channel
type BroadcastMessage struct {
	Channel string
	//Message []byte
	Data []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients by channel
	// Key: channel name (e.g., "agent::uuid" or "admin::dashboard")
	channels map[string]map[*Client]bool

	// Register requests from the clients
	register chan *Client

	router *MessageRouter

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast message to specific channel
	broadcast chan *BroadcastMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	hub := &Hub{
		broadcast:  make(chan *BroadcastMessage), // ✅ Правильный тип
		register:   make(chan *Client),
		unregister: make(chan *Client),
		channels:   make(map[string]map[*Client]bool),
	}

	// Initialize router
	hub.router = NewMessageRouter(hub)

	return hub
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.RegisterClient(client)

		case client := <-h.unregister:
			h.UnregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToChannel(message.Channel, message.Data)
		}
	}
}

func (h *Hub) SetRouter(r *MessageRouter) {
	h.router = r
}

// registerClient registers a new client to a channel
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.channels[client.Channel] == nil {
		h.channels[client.Channel] = make(map[*Client]bool)
	}
	h.channels[client.Channel][client] = true

	logrus.WithFields(logrus.Fields{
		"channel":   client.Channel,
		"client_id": client.ID,
		"user_type": client.UserType,
	}).Info("Client registered to channel")
}

// unregisterClient unregisters a client from its channel
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.channels[client.Channel]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)

			// Remove channel if empty
			if len(clients) == 0 {
				delete(h.channels, client.Channel)
			}

			logrus.WithFields(logrus.Fields{
				"channel":   client.Channel,
				"client_id": client.ID,
				"user_type": client.UserType,
			}).Info("Client unregistered from channel")
		}
	}
}

// broadcastToChannel sends a message to all clients in a specific channel
func (h *Hub) broadcastToChannel(channel string, data []byte) {
	h.mu.RLock()
	clients, exists := h.channels[channel]
	h.mu.RUnlock()

	if !exists {
		logrus.WithField("channel", channel).Debug("Channel not found for broadcast")
		return
	}

	for client := range clients {
		select {
		case client.send <- data:
		default:
			logrus.WithFields(logrus.Fields{
				"client_id": client.ID,
				"channel":   channel,
			}).Warn("Client send buffer full, skipping message")
		}
	}
}

// BroadcastToChannel broadcasts a message to all clients in a specific channel
func (h *Hub) BroadcastToChannel(channel string, msg *Message) {
	h.mu.RLock()
	clients, exists := h.channels[channel]
	h.mu.RUnlock()

	if !exists {
		logrus.WithField("channel", channel).Debug("Channel not found for broadcast")
		return
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal message for broadcast")
		return
	}

	for client := range clients {
		select {
		case client.send <- msgBytes:
		default:
			logrus.WithFields(logrus.Fields{
				"client_id": client.ID,
				"channel":   channel,
			}).Warn("Client send buffer full, skipping message")
		}
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
				"channel":        client.Channel,
				"client_id":      client.ID,
				"last_heartbeat": client.lastHeartbeat,
			}).Warn("Removing stale connection")
			h.unregister <- client
		}
	}
}

// SendToChannel sends a message to first available client in channel (for agent targeting)
func (h *Hub) SendToChannel(channel string, msg *Message) {
	h.mu.RLock()
	clients, exists := h.channels[channel]
	h.mu.RUnlock()

	if !exists {
		logrus.WithField("channel", channel).Warn("Channel not found for message")
		return
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal message")
		return
	}

	// Send to first available client (typically only one agent per channel)
	for client := range clients {
		select {
		case client.send <- msgBytes:
			logrus.WithFields(logrus.Fields{
				"client_id": client.ID,
				"channel":   channel,
			}).Debug("Message sent to client")
			return
		default:
			logrus.WithFields(logrus.Fields{
				"client_id": client.ID,
				"channel":   channel,
			}).Warn("Client send buffer full")
		}
	}
}
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}
