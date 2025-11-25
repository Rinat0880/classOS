package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 30 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024 // 512 KB
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// Unique client ID (user_id for admins, agent_id for agents)
	ID string

	// User type: "admin" or "agent"
	UserType string

	// Channel name this client is subscribed to
	Channel string

	// The hub instance
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Last heartbeat timestamp
	lastHeartbeat time.Time
}

// NewClient creates a new WebSocket client
func NewClient(id string, userType string, channel string, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		ID:            id,
		UserType:      userType,
		Channel:       channel,
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, 256),
		lastHeartbeat: time.Now(),
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c // Используем канал unregister
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Get router from hub
	router := c.hub.router
	if router == nil {
		logrus.Error("Message router not initialized")
		return
	}

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("WebSocket error")
			}
			break
		}

		// Route message through router
		if err := router.RouteMessage(c, message); err != nil {
			logrus.WithError(err).Error("Failed to route message")
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes different message types
func (c *Client) handleMessage(msg *Message) {
	logrus.WithFields(logrus.Fields{
		"type":      msg.Type,
		"client_id": c.ID,
		"channel":   c.Channel,
	}).Debug("Message received")

	switch msg.Type {
	case "status_update":
		// Handle status update from agent
		if c.UserType == "agent" {
			// Передаем msg напрямую, Hub сам сделает JSON marshal
			c.hub.BroadcastToChannel("admin::dashboard", msg)
		}

	case "action_log":
		// Handle action log from agent
		if c.UserType == "agent" {
			// Broadcast to admin dashboard
			c.hub.BroadcastToChannel("admin::dashboard", msg)
		}

	case "command_request":
		// Handle command request from admin to agent
		if c.UserType == "admin" {
			// !!! ВНИМАНИЕ: Тут есть еще одна ошибка (см. ниже)
			// Payload это []byte, его нельзя читать как map
			// Вам нужно сначала распаковать его

			// Временное решение для примера (нужна структура для распаковки):
			var params map[string]interface{}
			if err := json.Unmarshal(msg.Payload, &params); err == nil {
				if targetChannel, ok := params["target_channel"].(string); ok {
					c.hub.BroadcastToChannel(targetChannel, msg)
				}
			}
		}

	case "command_response":
		// Handle command response from agent
		if c.UserType == "agent" {
			// Forward to admin dashboard
			c.hub.BroadcastToChannel("admin::dashboard", msg)
		}

	default:
		logrus.WithField("type", msg.Type).Warn("Unknown message type")
	}
}

// marshalMessage converts a message to JSON bytes
func (c *Client) marshalMessage(msg *Message) []byte {
	data, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal message")
		return []byte{}
	}
	return data
}
