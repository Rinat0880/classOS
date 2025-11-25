package websocket

import (
	"encoding/json"
	"time"
)

// Message types
const (
	// Agent -> Backend
	MessageTypeHeartbeat       = "heartbeat"
	MessageTypeStatusUpdate    = "status_update"
	MessageTypeActionLog       = "action_log"
	MessageTypeCommandResponse = "command_response"

	// Admin -> Backend -> Agent
	MessageTypeCommandRequest = "command_request"

	// Backend -> Admin
	MessageTypeAgentStatus = "agent_status"
)

// Message represents a WebSocket message
type Message struct {
	Type      string          `json:"type"`                 // Message type
	Payload   json.RawMessage `json:"payload"`              // Flexible payload
	Timestamp time.Time       `json:"timestamp"`            // Message timestamp
	AgentID   string          `json:"agent_id,omitempty"`   // Agent ID (if applicable)
	RequestID string          `json:"request_id,omitempty"` // For request/response correlation
}

// HeartbeatPayload represents heartbeat message payload
type HeartbeatPayload struct {
	Status      string  `json:"status"`       // online, idle, busy
	CPUUsage    float64 `json:"cpu_usage"`    // CPU usage percentage
	MemoryUsage float64 `json:"memory_usage"` // Memory usage percentage
	Username    string  `json:"username"`     // Current logged in user
}

// StatusUpdatePayload represents status change message
type StatusUpdatePayload struct {
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
	Reason    string `json:"reason,omitempty"`
}

// ActionLogPayload represents agent action/event
type ActionLogPayload struct {
	Action   string                 `json:"action"`             // process_blocked, url_blocked, file_accessed
	Resource string                 `json:"resource"`           // Process name, URL, file path
	Username string                 `json:"username"`           // User who triggered action
	Allowed  bool                   `json:"allowed"`            // Was action allowed
	Metadata map[string]interface{} `json:"metadata,omitempty"` // Additional context
}

// CommandRequestPayload represents command from admin to agent
type CommandRequestPayload struct {
	Command string                 `json:"command"` // get_processes, kill_process, update_whitelist
	Params  map[string]interface{} `json:"params"`  // Command parameters
	Timeout int                    `json:"timeout"` // Timeout in seconds
}

// CommandResponsePayload represents agent's response to command
type CommandResponsePayload struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// AgentStatusPayload for broadcasting agent status to admins
type AgentStatusPayload struct {
	AgentID  string    `json:"agent_id"`
	Hostname string    `json:"hostname"`
	Status   string    `json:"status"` // online, offline
	LastSeen time.Time `json:"last_seen"`
	Username string    `json:"username,omitempty"`
}

// NewMessage creates a new message with timestamp
func NewMessage(msgType string, payload interface{}, agentID string) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:      msgType,
		Payload:   payloadBytes,
		Timestamp: time.Now(),
		AgentID:   agentID,
	}, nil
}

// DecodePayload decodes message payload into target struct
func (m *Message) DecodePayload(target interface{}) error {
	return json.Unmarshal(m.Payload, target)
}
