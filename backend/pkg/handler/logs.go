package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	classosbackend "github.com/rinat0880/classOS_backend"
)

func (h *Handler) getLogs(c *gin.Context) {
	username := c.Query("username")
	deviceName := c.Query("device")
	fromDateStr := c.Query("from")
	toDateStr := c.Query("to")
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	filter := classosbackend.LogsFilter{
		Username:   username,
		DeviceName: deviceName,
		Limit:      limit,
		Offset:     offset,
	}

	if fromDateStr != "" {
		fromDate, err := time.Parse(time.RFC3339, fromDateStr)
		if err == nil {
			filter.FromDate = &fromDate
		}
	}

	if toDateStr != "" {
		toDate, err := time.Parse(time.RFC3339, toDateStr)
		if err == nil {
			filter.ToDate = &toDate
		}
	}

	logs, err := h.services.Logs.GetLogsFiltered(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	count, err := h.services.Logs.GetLogsCount(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data":  logs,
		"total": count,
		"limit": limit,
		"offset": offset,
	})
}

func (h *Handler) getLogsByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		newErrorResponse(c, http.StatusBadRequest, "username is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	logs, err := h.services.Logs.GetLogsByUsername(username, limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": logs,
	})
}

func (h *Handler) getLogsByDevice(c *gin.Context) {
	deviceName := c.Param("device")
	if deviceName == "" {
		newErrorResponse(c, http.StatusBadRequest, "device name is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	logs, err := h.services.Logs.GetLogsByDevice(deviceName, limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": logs,
	})
}
