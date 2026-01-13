package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllDevices(c *gin.Context) {
	devices, err := h.services.Device.GetAllDevices()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": devices,
	})
}

func (h *Handler) getOnlineDevices(c *gin.Context) {
	devices, err := h.services.Device.GetOnlineDevices()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": devices,
	})
}

func (h *Handler) getDeviceByName(c *gin.Context) {
	deviceName := c.Param("name")
	if deviceName == "" {
		newErrorResponse(c, http.StatusBadRequest, "device name is required")
		return
	}

	device, err := h.services.Device.GetDeviceByName(deviceName)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "device not found")
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *Handler) deleteDevice(c *gin.Context) {
	deviceName := c.Param("name")
	if deviceName == "" {
		newErrorResponse(c, http.StatusBadRequest, "device name is required")
		return
	}

	err := h.services.Device.DeleteDevice(deviceName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
