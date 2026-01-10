package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "checkerId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")

	// ВРЕМЕННОЕ ЛОГИРОВАНИЕ
	logrus.Info("=== userIdentity middleware ===")
	logrus.WithField("header", header).Info("Authorization header")
	logrus.WithField("method", c.Request.Method).Info("Request method")
	logrus.WithField("path", c.Request.URL.Path).Info("Request path")

	if header == "" {
		logrus.Error("empty auth header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		logrus.Error("invalid auth header format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}

	checkerId, role, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		logrus.WithError(err).Error("token parse failed")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logrus.WithFields(logrus.Fields{
		"checkerId": checkerId,
		"role":      role,
	}).Info("authentication successful")
	c.Set("checkerId", checkerId)
	c.Set("role", role)
	c.Next()
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "invalid type of id")
		return 0, errors.New("user id not found")
	}

	return idInt, nil
}

func (h *Handler) adminOnly(c *gin.Context) {
	roleVal, exists := c.Get("role")
	if !exists {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found"})
		return
	}

	role := roleVal.(string)
	if role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	c.Next()
}
