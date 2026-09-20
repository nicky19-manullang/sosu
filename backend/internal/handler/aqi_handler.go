package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sosu-backend/internal/repository"
)

type AQIHandler struct {
	repo *repository.AQIRepository
}

func NewAQIHandler(repo *repository.AQIRepository) *AQIHandler {
	return &AQIHandler{repo: repo}
}

func (h *AQIHandler) GetNearest(c *gin.Context) {
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

	reading, distanceMeters, err := h.repo.FindNearest(lat, lng, 75000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if reading == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tidak ada data AQI dalam radius 75km"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reading":       reading,
		"distance_km":   distanceMeters / 1000,
	})
}