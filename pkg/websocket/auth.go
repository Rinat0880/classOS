package websocket

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware handles WebSocket authorization
type AuthMiddleware struct {
	signingKey []byte
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(signingKey string) *AuthMiddleware {
	return &AuthMiddleware{
		signingKey: []byte(signingKey),
	}
}

// AuthorizeWebSocket validates token and returns client info
func (am *AuthMiddleware) AuthorizeWebSocket(c *gin.Context) (clientID string, userType string, channel string, err error) {
	// Get token from query parameter or header
	token := c.Query("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		return "", "", "", errors.New("missing authorization token")
	}

	// Determine if it's JWT (admin) or device token (agent)
	if strings.HasPrefix(token, "device_") {
		// Device token for agents
		return am.validateDeviceToken(token)
	}

	// JWT token for admins
	return am.validateJWT(token)
}

// validateJWT validates JWT token for admin users
func (am *AuthMiddleware) validateJWT(tokenString string) (string, string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return am.signingKey, nil
	})

	if err != nil {
		logrus.WithError(err).Warn("Failed to parse JWT token")
		return "", "", "", errors.New("invalid token")
	}

	if !token.Valid {
		return "", "", "", errors.New("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", errors.New("invalid token claims")
	}

	// Extract checker_id (user_id) and role
	checkerID, ok := claims["checker_id"].(float64)
	if !ok {
		return "", "", "", errors.New("invalid checker_id in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", "", errors.New("invalid role in token")
	}

	if role != "admin" {
		return "", "", "", errors.New("only admins can connect via JWT")
	}

	clientID := fmt.Sprintf("%d", int(checkerID))
	channel := "admin::dashboard"

	logrus.WithFields(logrus.Fields{
		"client_id": clientID,
		"role":      role,
		"channel":   channel,
	}).Info("Admin authenticated via JWT")

	return clientID, "admin", channel, nil
}

// validateDeviceToken validates device token for agent connections
func (am *AuthMiddleware) validateDeviceToken(token string) (string, string, string, error) {
	// TODO: Query database for device token validation
	// For now, basic validation

	if !strings.HasPrefix(token, "device_") {
		return "", "", "", errors.New("invalid device token format")
	}

	// Extract agent_id from token (simplified for now)
	agentID := strings.TrimPrefix(token, "device_")
	if agentID == "" {
		return "", "", "", errors.New("empty agent_id in device token")
	}

	channel := fmt.Sprintf("agent::%s", agentID)

	logrus.WithFields(logrus.Fields{
		"agent_id": agentID,
		"channel":  channel,
	}).Info("Agent authenticated via device token")

	return agentID, "agent", channel, nil
}

// WSAuthMiddleware is a Gin middleware for WebSocket authorization
func (am *AuthMiddleware) WSAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, userType, channel, err := am.AuthorizeWebSocket(c)
		if err != nil {
			logrus.WithError(err).Warn("WebSocket authorization failed")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Store in context for handler
		c.Set("ws_client_id", clientID)
		c.Set("ws_user_type", userType)
		c.Set("ws_channel", channel)

		c.Next()
	}
}
