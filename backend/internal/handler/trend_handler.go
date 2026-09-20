package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sosu-backend/internal/service"
)

type TrendHandler struct {
	trendService *service.TrendService
}

func NewTrendHandler(trendService *service.TrendService) *TrendHandler {
	return &TrendHandler{trendService: trendService}
}

func (h *TrendHandler) GetTrend(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter lat wajib diisi"})
		return
	}
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter lng wajib diisi"})
		return
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 30 {
			days = parsed
		}
	}

	points, err := h.trendService.GetTrend(lat, lng, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trend": points})
}