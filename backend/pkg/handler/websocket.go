package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	classosbackend "github.com/rinat0880/classOS_backend"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Type      string                    `json:"type"`
	Device    string                    `json:"device,omitempty"`
	User      string                    `json:"user,omitempty"`
	Timestamp time.Time                 `json:"timestamp,omitempty"`
	Data      []classosbackend.UserLog  `json:"data,omitempty"`
	Token     string                    `json:"token,omitempty"`
}

func (h *Handler) handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	authenticated := false
	var deviceName string

	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		switch msg.Type {
		case "auth":
			if msg.Token == "" {
				log.Printf("Authentication failed: no token")
				conn.Close()
				return
			}
			authenticated = true
			log.Printf("Agent authenticated with token")

		case "heartbeat":
			if !authenticated {
				log.Printf("Heartbeat from unauthenticated client")
				continue
			}

			deviceName = msg.Device
			device := classosbackend.DeviceStatus{
				DeviceName:    msg.Device,
				Username:      msg.User,
				LastHeartbeat: time.Now(),
			}

			err := h.services.Device.UpsertDeviceStatus(device)
			if err != nil {
				log.Printf("Failed to update device status: %v", err)
			} else {
				log.Printf("Heartbeat received from device %s, user %s", msg.Device, msg.User)
			}

		case "logs":
			if !authenticated {
				log.Printf("Logs from unauthenticated client")
				continue
			}

			if len(msg.Data) > 0 {
				for i := range msg.Data {
					if msg.Data[i].DeviceName == "" && deviceName != "" {
						msg.Data[i].DeviceName = deviceName
					}
				}

				err := h.services.Logs.SaveLogs(msg.Data)
				if err != nil {
					log.Printf("Failed to save logs: %v", err)
				} else {
					log.Printf("Saved %d logs from device %s", len(msg.Data), deviceName)
				}
			}

		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}

		response := map[string]string{"status": "ok"}
		jsonResp, _ := json.Marshal(response)
		conn.WriteMessage(websocket.TextMessage, jsonResp)
	}
}
