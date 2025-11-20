package websocket

import "time"

// Message represents a WebSocket message structure
type Message struct {
	// Message type: "heartbeat", "status_update", "action_log", "command_request", "command_response"
	Type string `json:"type"`

	// Message payload (structure depends on type)
	Payload map[string]interface{} `json:"payload"`

	// Timestamp when message was created
	Timestamp time.Time `json:"timestamp"`

	// Client ID (agent_id or admin_id)
	ClientID string `json:"client_id,omitempty"`
}

// HeartbeatPayload represents the payload for heartbeat messages
type HeartbeatPayload struct {
	AgentID      string `json:"agent_id"`
	Hostname     string `json:"hostname"`
	CurrentUser  string `json:"current_user"`
	AgentVersion string `json:"agent_version"`
}

// StatusUpdatePayload represents the payload for status update messages
type StatusUpdatePayload struct {
	AgentID string `json:"agent_id"`
	Status  string `json:"status"` // "online", "offline", "busy"
	Details string `json:"details,omitempty"`
}

// ActionLogPayload represents the payload for action log messages
type ActionLogPayload struct {
	EventType string                 `json:"event_type"` // "program_launched", "file_opened", "project_saved"
	Process   string                 `json:"process,omitempty"`
	Path      string                 `json:"path,omitempty"`
	Username  string                 `json:"username"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

// CommandRequestPayload represents the payload for command request messages
type CommandRequestPayload struct {
	Command       string                 `json:"command"` // "get_screenshot", "restart", etc.
	TargetChannel string                 `json:"target_channel"`
	RequestID     string                 `json:"request_id"`
	Parameters    map[string]interface{} `json:"parameters,omitempty"`
}

// CommandResponsePayload represents the payload for command response messages
type CommandResponsePayload struct {
	RequestID string                 `json:"request_id"`
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
}
