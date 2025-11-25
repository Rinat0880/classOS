package websocket

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// AgentWSClient represents a WebSocket client for agent with reconnection logic
type AgentWSClient struct {
	serverURL          string
	deviceToken        string
	conn               *websocket.Conn
	mu                 sync.Mutex
	reconnectAttempts  int
	maxReconnectDelay  time.Duration
	baseReconnectDelay time.Duration
	done               chan struct{}
	messageHandlers    map[string]func(*Message) error
}

// NewAgentWSClient creates a new agent WebSocket client
func NewAgentWSClient(serverURL, deviceToken string) *AgentWSClient {
	return &AgentWSClient{
		serverURL:          serverURL,
		deviceToken:        deviceToken,
		maxReconnectDelay:  30 * time.Second,
		baseReconnectDelay: 1 * time.Second,
		done:               make(chan struct{}),
		messageHandlers:    make(map[string]func(*Message) error),
	}
}

// Connect establishes WebSocket connection
func (c *AgentWSClient) Connect() error {
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}

	// Add device token as query parameter
	q := u.Query()
	q.Set("token", c.deviceToken)
	u.RawQuery = q.Encode()

	logrus.WithField("url", u.String()).Info("Connecting to WebSocket server")

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.reconnectAttempts = 0 // Reset on successful connection
	c.mu.Unlock()

	logrus.Info("WebSocket connection established")

	// Start read/write loops
	go c.readLoop()
	go c.heartbeatLoop()

	return nil
}

// ConnectWithRetry attempts to connect with exponential backoff
func (c *AgentWSClient) ConnectWithRetry() {
	for {
		err := c.Connect()
		if err == nil {
			return
		}

		// Calculate exponential backoff delay
		delay := time.Duration(math.Min(
			float64(c.baseReconnectDelay)*math.Pow(2, float64(c.reconnectAttempts)),
			float64(c.maxReconnectDelay),
		))

		c.reconnectAttempts++
		logrus.WithFields(logrus.Fields{
			"attempt": c.reconnectAttempts,
			"delay":   delay,
			"error":   err,
		}).Warn("Failed to connect, retrying...")

		select {
		case <-time.After(delay):
			continue
		case <-c.done:
			return
		}
	}
}

// readLoop continuously reads messages from WebSocket
func (c *AgentWSClient) readLoop() {
	defer c.scheduleReconnect()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logrus.Info("WebSocket connection closed normally")
			} else {
				logrus.WithError(err).Error("Error reading message")
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			logrus.WithError(err).Error("Failed to unmarshal message")
			continue
		}

		// Handle message based on type
		if handler, exists := c.messageHandlers[msg.Type]; exists {
			if err := handler(&msg); err != nil {
				logrus.WithError(err).WithField("type", msg.Type).Error("Error handling message")
			}
		} else {
			logrus.WithField("type", msg.Type).Warn("No handler for message type")
		}
	}
}

// heartbeatLoop sends periodic heartbeats
func (c *AgentWSClient) heartbeatLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.SendHeartbeat(); err != nil {
				logrus.WithError(err).Error("Failed to send heartbeat")
			}
		case <-c.done:
			return
		}
	}
}

// SendHeartbeat sends heartbeat message to server
func (c *AgentWSClient) SendHeartbeat() error {
	payload := HeartbeatPayload{
		Status:      "online",
		CPUUsage:    0.0, // TODO: Get actual CPU usage
		MemoryUsage: 0.0, // TODO: Get actual memory usage
		Username:    "",  // TODO: Get current logged in user
	}

	return c.Send(MessageTypeHeartbeat, payload, "")
}

// Send sends a message through WebSocket
func (c *AgentWSClient) Send(msgType string, payload interface{}, agentID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("WebSocket connection not established")
	}

	msg, err := NewMessage(msgType, payload, agentID)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// On registers a message handler for specific message type
func (c *AgentWSClient) On(msgType string, handler func(*Message) error) {
	c.messageHandlers[msgType] = handler
}

// scheduleReconnect attempts to reconnect after connection loss
func (c *AgentWSClient) scheduleReconnect() {
	logrus.Info("Scheduling reconnection...")
	go c.ConnectWithRetry()
}

// Close closes the WebSocket connection
func (c *AgentWSClient) Close() {
	close(c.done)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
	}
}
