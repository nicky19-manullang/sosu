package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sosu-backend/internal/service"
)

type StatusHandler struct {
	statusService *service.StatusService
}

func NewStatusHandler(statusService *service.StatusService) *StatusHandler {
	return &StatusHandler{statusService: statusService}
}

func (h *StatusHandler) GetStatus(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter lat wajib diisi dan berupa angka"})
		return
	}
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter lng wajib diisi dan berupa angka"})
		return
	}

	result, err := h.statusService.GetStatus(lat, lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}