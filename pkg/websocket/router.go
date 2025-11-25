package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// MessageRouter handles routing of messages between channels
type MessageRouter struct {
	hub *Hub
}

// NewMessageRouter creates a new message router
func NewMessageRouter(hub *Hub) *MessageRouter {
	return &MessageRouter{
		hub: hub,
	}
}

// RouteMessage routes incoming message to appropriate handlers
func (r *MessageRouter) RouteMessage(client *Client, data []byte) error {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal message")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"type":      msg.Type,
		"client_id": client.ID,
		"agent_id":  msg.AgentID,
		"channel":   client.Channel,
	}).Debug("Routing message")

	switch msg.Type {
	case MessageTypeHeartbeat:
		return r.handleHeartbeat(client, &msg)
	case MessageTypeStatusUpdate:
		return r.handleStatusUpdate(client, &msg)
	case MessageTypeActionLog:
		return r.handleActionLog(client, &msg)
	case MessageTypeCommandRequest:
		return r.handleCommandRequest(client, &msg)
	case MessageTypeCommandResponse:
		return r.handleCommandResponse(client, &msg)
	default:
		logrus.WithField("type", msg.Type).Warn("Unknown message type")
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

// handleHeartbeat processes heartbeat from agent
func (r *MessageRouter) handleHeartbeat(client *Client, msg *Message) error {
	var payload HeartbeatPayload
	if err := msg.DecodePayload(&payload); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"agent_id": client.ID,
		"status":   payload.Status,
		"username": payload.Username,
		"cpu":      payload.CPUUsage,
		"memory":   payload.MemoryUsage,
	}).Debug("Heartbeat received")

	// Broadcast agent status to admin dashboard
	statusPayload := AgentStatusPayload{
		AgentID:  client.ID,
		Hostname: client.ID, // TODO: Get actual hostname from agent registration
		Status:   "online",
		LastSeen: msg.Timestamp,
		Username: payload.Username,
	}

	statusMsg, err := NewMessage(MessageTypeAgentStatus, statusPayload, client.ID)
	if err != nil {
		return err
	}

	// Broadcast to admin::dashboard channel
	r.hub.BroadcastToChannel("admin::dashboard", statusMsg)

	return nil
}

// handleStatusUpdate processes status change from agent
func (r *MessageRouter) handleStatusUpdate(client *Client, msg *Message) error {
	var payload StatusUpdatePayload
	if err := msg.DecodePayload(&payload); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"agent_id":   client.ID,
		"old_status": payload.OldStatus,
		"new_status": payload.NewStatus,
		"reason":     payload.Reason,
	}).Info("Status update received")

	// Broadcast to admin dashboard
	r.hub.BroadcastToChannel("admin::dashboard", msg)

	return nil
}

// handleActionLog processes action log from agent
func (r *MessageRouter) handleActionLog(client *Client, msg *Message) error {
	var payload ActionLogPayload
	if err := msg.DecodePayload(&payload); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"agent_id": client.ID,
		"action":   payload.Action,
		"resource": payload.Resource,
		"username": payload.Username,
		"allowed":  payload.Allowed,
	}).Info("Action log received")

	// TODO: Store in database (agent_logs table)

	// Broadcast to admin dashboard
	r.hub.BroadcastToChannel("admin::dashboard", msg)

	return nil
}

// handleCommandRequest processes command from admin to agent
func (r *MessageRouter) handleCommandRequest(client *Client, msg *Message) error {
	if client.UserType != "admin" {
		return fmt.Errorf("only admins can send commands")
	}

	var payload CommandRequestPayload
	if err := msg.DecodePayload(&payload); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"admin_id": client.ID,
		"agent_id": msg.AgentID,
		"command":  payload.Command,
	}).Info("Command request received")

	// Route to specific agent channel
	targetChannel := fmt.Sprintf("agent::%s", msg.AgentID)
	r.hub.SendToChannel(targetChannel, msg)

	return nil
}

// handleCommandResponse processes response from agent to command
func (r *MessageRouter) handleCommandResponse(client *Client, msg *Message) error {
	var payload CommandResponsePayload
	if err := msg.DecodePayload(&payload); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"agent_id":   client.ID,
		"request_id": msg.RequestID,
		"success":    payload.Success,
	}).Info("Command response received")

	// TODO: Route back to specific admin who sent the command (need to track request_id)
	// For now, broadcast to all admins
	r.hub.BroadcastToChannel("admin::dashboard", msg)

	return nil
}
